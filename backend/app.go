package backend

import (
	"HLAE-Studio/backend/api"
	"HLAE-Studio/backend/config"
	"HLAE-Studio/backend/tool"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/otiai10/copy"
	"github.com/wailsapp/wails"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

///// app.go 存放backend包与frontend交互的大部分操作

//App设置
type App struct {
	runtime *wails.Runtime //初始化Runtime需要
	cfg     config.CFG
}

//Wails初始化
func (a *App) WailsInit(runtime *wails.Runtime) error {
	a.runtime = runtime
	var err error
	if a.cfg, err = config.ReadConfig("./config.json"); err != nil {
		a.runtime.Events.Emit("SetLog", err)
		log.Println(err)
		return err
	}


	return nil
}

//Wails结束前
func (a *App) WailsShutdown() {
	//结束前：
	fmt.Println("Wails结束")
	err := config.SaveConfig(a.cfg, "./config.json")
	if err != nil {
		log.Println(err)
	}
}

//设置前端变量
func (a *App) SetVar() {
	a.setAppVersion(a.cfg.AppVersion)
	a.setVersionCode(a.cfg.VersionCode)
}

//检查更新
func (a *App) CheckUpdate() {
	var err error
	//我更新我自己
	if err = a.updateApp(); err != nil {
		a.noticeError(err.Error())
	}

	//检查是否初始化，否则通知前端选择安装方式和位置
	if a.cfg.Init == false {
		a.doSelectOption()
		return
	}

	//检查HLAE和FFmpeg安装情况
	if err = a.checkState(); err != nil {
		a.runtime.Events.Emit("SetLog", err)
		log.Println(err)
		return
	}

	//安装或者更新HLAE
	if a.cfg.HlaeState == false {
		err = a.installHLAE()
	} else {
		err = a.updateHLAE()
	}
	if err != nil {
		log.Println(err)
		a.setLog(err.Error())
		return
	}

	//安装或更新FFmpeg
	if a.cfg.FFmpegState == false {
		err = a.installFFmpeg()
	} else {
		err = a.updateFFmpeg()
	}
	if err != nil {
		log.Println(err)
		a.setLog(err.Error())
		return
	}

	a.setProgress(100)
	a.setLog("当前是最新版本")
	a.noticeSuccess("已经更新到最新版本")
	a.cfg.Init = true
	a.cfg.HlaeState = true
	a.cfg.FFmpegState = true
	return
}

//检查HLAE和FFmpeg安装情况
func (a *App) checkState() error {
	//如果已经初始化，再核实一下hlae和ffmpeg程序位置
	if a.cfg.Init == true {
		//检查hlae
		if ok, err := tool.IsFileExisted(a.cfg.HlaePath + "/hlae.exe"); err != nil {
			return err
		} else if ok == true {
			a.cfg.HlaeState = true
			//解析修正本地hlae版本号
			changelog, err := tool.ReadAll(a.cfg.HlaePath + "/changelog.xml")
			if err != nil {
				a.noticeError("读取本地版本号失败: " + err.Error())
				return err
			}
			if tVersion, err := api.ParseChangelog(changelog); err != nil {
				a.noticeError("解析本地版本号失败: " + err.Error())
				return err
			} else {
				a.cfg.HlaeVersion = tVersion
			}
			//检查ffmpeg
			if ok, err := tool.IsFileExisted(a.cfg.HlaePath + "/ffmpeg/bin/ffmpeg.exe"); err != nil {
				return err
			} else if ok == true {
				a.cfg.FFmpegState = true
			} else {
				a.cfg.FFmpegState = false
				a.cfg.Init = false
			}
		} else {
			a.cfg.HlaeState = false
			a.cfg.FFmpegState = false
			a.cfg.HlaePath = ""
			a.cfg.Init = false
		}
	} else {
		//否则直接给定false
		a.cfg.HlaeState = false
		a.cfg.FFmpegState = false
		a.cfg.HlaePath = ""
		a.cfg.Init = false
	}

	return nil
}

//TODO 我更新我自己
func (a *App) updateApp() error {

	return nil
}

//设定安装hlae的方式
func (a *App) SetOption(ok bool) {
	a.cfg.Standalone = ok
	a.cfg.Init = true
}

//安装HLAE
func (a *App) installHLAE() error {
	//选择安装位置
	if a.cfg.Standalone == true {
		a.noticeWarning("请选择安装位置，或者已有hlae.exe的文件夹")
		//time.Sleep(2 * time.Second)
		path := a.SelectDirectory()
		if path == "" {
			a.noticeWarning("已取消安装")
			_ = tool.WriteFast("./cancel.txt", "取消安装")
			return nil
		}

		//识别已安装hlae，则检查更新
		if ok, _ := tool.IsFileExisted(path + "/hlae.exe"); ok {
			a.cfg.HlaePath = tool.FormatPath(path)
			//解析修正本地hlae版本号
			changelog, err := tool.ReadAll(a.cfg.HlaePath + "/changelog.xml")
			if err != nil {
				a.noticeError("读取本地版本号失败: " + err.Error())
				return err
			}
			if tVersion, err := api.ParseChangelog(changelog); err != nil {
				a.noticeError("解析本地版本号失败: " + err.Error())
				return err
			} else {
				a.cfg.HlaeVersion = tVersion
			}
			if err := a.updateHLAE(); err != nil {
				return err
			}
			return nil
		}

		//新装hlae
		a.cfg.HlaePath = tool.FormatPath(path)

	} else {
		//检查CSGO Demos Manager是否安装 "%HOMEDIR%/AppData/Local/AkiVer/"
		usr, err := user.Current()
		if err != nil {
			return err
		}
		if ok, err := tool.IsFileExisted(usr.HomeDir + "/AppData/Local/AkiVer/"); err != nil || !ok {
			//Manager未安装  提示安装且弹出下载的网页
			if err := a.runtime.Browser.OpenURL("https://github.com/akiver/CSGO-Demos-Manager/releases/latest"); err != nil {
				return err
			}
			if err := a.runtime.Browser.OpenURL("https://hlae.site/topic/390"); err != nil {
				return err
			}

			return errors.New("CSGO Demos Manager未安装")
		} else {
			a.cfg.HlaePath = tool.FormatPath(usr.HomeDir + "/AppData/Local/AkiVer/hlae")
		}
	}

	//声明变量
	var srcVersion, cdnVersion, srcURL, cdnURL, srcFilename, cdnFilename, version, url, filename string

	//声明用来控制进度条的变量
	const count = 5
	i := 0
	a.setProgress(0)
	a.setLog("正在读取HLAE和CDN源API...")
	//尝试读取HLAE和CDN源API
	srcData, err1 := tool.GetHttpData(a.cfg.HlaeAPI)
	cdnData, err2 := tool.GetHttpData(a.cfg.HlaeCdnAPI)
	if err1 != nil && err2 != nil {
		return errors.New("hlae官方和CDN的API均获取失败")
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在解析API...")
	//解析API信息，决定下载的文件
	if err1 == nil {
		//解析官方API
		var latestInst api.GitHubLatest
		//注释下面一行->使用encoding/json库
		var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary //使用高性能json-iterator/go库
		err := jsonx.Unmarshal([]byte(srcData), &latestInst)      //第二个参数要地址传递
		if err != nil || latestInst.Message == "Not Found" || strings.Contains(latestInst.Message, "API rate limit"){
			//HLAE API获取失败
			a.noticeWarning("官方API解析失败: " + latestInst.Message)
		} else {
			srcVersion = latestInst.TagName
			//获得zip附件信息
			for _, file := range latestInst.Assets {
				//过滤掉源码文件
				if file.State == "uploaded" && !strings.Contains(file.Name, ".asc") && strings.Contains(file.Name, ".zip") {
					srcURL = file.BrowserDownloadURL
					srcFilename = file.Name
				}
			}
		}
	}

	if err2 == nil {
		//解析CDN API
		var cdnInst api.ReleaseDelivr
		//注释下面一行->使用encoding/json库
		var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary //使用高性能json-iterator/go库
		err := jsonx.Unmarshal([]byte(cdnData), &cdnInst)      //第二个参数要地址传递
		if err != nil {
			return err
		}
		//获取版本号、下载地址和文件名
		cdnVersion = cdnInst.Version
		cdnURL = cdnInst.DownloadLink[0]
		_, cdnFilename = filepath.Split(cdnURL)
	}

	//决定下载的文件
	if srcVersion == "" && cdnVersion == "" {
		return errors.New("hlae官方和CDN的API均获取或解析失败")
	} else if cdnVersion == "" {
		a.noticeWarning("CDN源解析失败，下载速度可能较慢")
	} else if srcVersion == "" {
		a.noticeWarning("hlae官方API解析失败，CDN源可能不是最新版本")
	}

	if srcVersion == "" || (srcVersion != "" && srcVersion == cdnVersion) {
		//官方版本非空且和CDN源版本一致，则下载CDN源
		version = cdnVersion
		url = cdnURL
		filename = cdnFilename
	} else {
		//否则下载官方源
		version = srcVersion
		url = srcURL
		filename = srcFilename
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在下载HLAE安装包...")
	//下载HLAE安装包
	if err := tool.DownloadFile(url, "./temp"); err != nil {
		a.noticeError("HLAE下载失败")
		return err
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在解压HLAE安装包...")
	//解压
	if err := tool.Decompress("./temp/" + filename, "./temp/hlae"); err != nil {
		return err
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在转移文件...")
	//转移
	if err := copy.Copy("./temp/hlae", a.cfg.HlaePath); err != nil {
		return err
	}
	_ = os.RemoveAll("./temp")

	//生成version文件
	a.cfg.HlaeVersion = version
	ver := strings.Replace(version, "v", "", -1)
	if err := tool.WriteFast(a.cfg.HlaePath + "/version", ver); err != nil {
		return err
	}

	//完成
	a.setProgress(100)
	a.setLog("HLAE安装完成")
	return nil
}

//安装FFmpeg
func (a *App) installFFmpeg() error {
	//声明变量
	var srcVersion, cdnVersion, srcURL, cdnURL, srcFilename, cdnFilename, version, url, filename string

	//声明用来控制进度条的变量
	const count = 5
	i := 0
	a.setProgress(0)
	a.setLog("正在读取FFmpeg和CDN源API...")
	//尝试读取FFmpeg和CDN源API
	srcData, err1 := tool.GetHttpData(a.cfg.FFmpegAPI + "/release-version")
	cdnData, err2 := tool.GetHttpData(a.cfg.FFmpegCdnAPI)
	if err1 != nil && err2 != nil {
		return errors.New("FFmpeg官方和CDN的API均获取失败")
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在解析API...")
	//解析API信息，决定下载的文件
	if err1 == nil {
		//解析官方API
		srcVersion = srcData
		srcURL = a.cfg.FFmpegAPI + "/ffmpeg-release-essentials.7z"
		srcFilename = "ffmpeg-release-essentials.7z"
	}

	if err2 == nil {
		//解析CDN API
		var cdnInst api.ReleaseDelivr
		//注释下面一行->使用encoding/json库
		var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary //使用高性能json-iterator/go库
		err := jsonx.Unmarshal([]byte(cdnData), &cdnInst)      //第二个参数要地址传递
		if err != nil {
			return err
		}
		//获取版本号、下载地址和文件名
		cdnVersion = cdnInst.Version
		cdnURL = cdnInst.DownloadLink[0]
		_, cdnFilename = filepath.Split(cdnURL)
	}

	//决定下载的文件
	if srcVersion == "" && cdnVersion == "" {
		return errors.New("FFmpeg官方和CDN的API均获取或解析失败")
	} else if srcVersion == "" {
		a.noticeWarning("FFmpeg官方API解析失败，CDN源可能不是最新版本")
	} else if cdnVersion == "" {
		a.noticeWarning("CDN源解析失败，下载速度可能较慢")
	}
	if srcVersion == "" || (srcVersion != "" && srcVersion == cdnVersion) {
		//官方版本非空且和CDN源版本一致，则下载CDN源
		version = cdnVersion
		url = cdnURL
		filename = cdnFilename
	} else {
		//否则下载官方源
		version = srcVersion
		url = srcURL
		filename = srcFilename
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在下载FFmpeg安装包...")
	//下载FFmpeg安装包
	if err := tool.DownloadFile(url, "./temp"); err != nil {
		a.noticeError("FFmpeg下载失败")
		return err
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在解压FFmpeg安装包...")
	//解压
	if err := tool.Decompress("./temp/" + filename, "./temp/ffmpeg"); err != nil {
		return err
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在转移文件...")

	//转移
	if ok, err := tool.IsFileExisted("./temp/ffmpeg/bin/ffmpeg.exe"); err != nil || !ok {
		if ok, err := tool.IsFileExisted("./temp/ffmpeg/ffmpeg-release-essentials/bin/ffmpeg.exe"); err != nil || !ok {
			return err
		} else {
			if err := copy.Copy("./temp/ffmpeg/ffmpeg-release-essentials", a.cfg.HlaePath + "/ffmpeg"); err != nil {
				return err	//\temp\ffmpeg\ffmpeg-release-essentials\bin
			}
		}
	} else {
		if err := copy.Copy("./temp/ffmpeg", a.cfg.HlaePath + "/ffmpeg"); err != nil {
			return err
		}
	}
	_ = os.RemoveAll("./temp")

	//完成
	a.cfg.FFmpegVersion = version
	a.setProgress(100)
	a.setLog("FFmpeg安装完成")
	return nil
}

//更新HLAE
func (a *App) updateHLAE() error {
	//声明变量
	var srcVersion, cdnVersion, srcURL, cdnURL, srcFilename, cdnFilename, version, url, filename string

	//声明用来控制进度条的变量
	const count = 5
	i := 0
	a.setProgress(0)
	a.setLog("正在读取HLAE和CDN源API...")
	//尝试读取HLAE和CDN源API
	srcData, err1 := tool.GetHttpData(a.cfg.HlaeAPI)
	cdnData, err2 := tool.GetHttpData(a.cfg.HlaeCdnAPI)
	if err1 != nil && err2 != nil {
		return errors.New("hlae官方和CDN的API均获取失败")
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在解析API...")
	//解析API信息，决定下载的文件
	if err1 == nil {
		//解析官方API
		var latestInst api.GitHubLatest
		//注释下面一行->使用encoding/json库
		var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary //使用高性能json-iterator/go库
		err := jsonx.Unmarshal([]byte(srcData), &latestInst)      //第二个参数要地址传递
		if err != nil || latestInst.Message == "Not Found" || strings.Contains(latestInst.Message, "API rate limit"){
			//HLAE API获取失败
			a.noticeWarning("官方API解析失败: " + latestInst.Message)
		} else {
			srcVersion = latestInst.TagName
			//获得zip附件信息
			for _, file := range latestInst.Assets {
				//过滤掉源码文件
				if file.State == "uploaded" && !strings.Contains(file.Name, ".asc") && strings.Contains(file.Name, ".zip") {
					srcURL = file.BrowserDownloadURL
					srcFilename = file.Name
				}
			}
		}
	}

	if err2 == nil {
		//解析CDN API
		var cdnInst api.ReleaseDelivr
		//注释下面一行->使用encoding/json库
		var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary //使用高性能json-iterator/go库
		err := jsonx.Unmarshal([]byte(cdnData), &cdnInst)      //第二个参数要地址传递
		if err != nil {
			return err
		}
		//获取版本号、下载地址和文件名
		cdnVersion = cdnInst.Version
		cdnURL = cdnInst.DownloadLink[0]
		_, cdnFilename = filepath.Split(cdnURL)
	}

	//决定是否更新和下载的文件
	if srcVersion == "" && cdnVersion == "" {
		return errors.New("hlae官方和CDN的API均获取或解析失败")
	} else if srcVersion == "" {
		a.noticeWarning("hlae官方API解析失败，CDN源可能不是最新版本")
	} else if cdnVersion == "" {
		a.noticeWarning("CDN源解析失败，下载速度可能较慢")
	}
	if srcVersion == "" || (srcVersion != "" && srcVersion == cdnVersion) {
		//官方版本非空且和CDN源版本一致，则下载CDN源
		version = cdnVersion
		url = cdnURL
		filename = cdnFilename
	} else {
		//否则下载官方源
		version = srcVersion
		url = srcURL
		filename = srcFilename
	}

	//对比版本号，决定是否更新
	if version == a.cfg.HlaeVersion {
		return nil
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在下载HLAE安装包...")
	//下载HLAE安装包
	if err := tool.DownloadFile(url, "./temp"); err != nil {
		a.noticeError("HLAE下载失败")
		return err
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在解压HLAE安装包...")
	//解压
	if err := tool.Decompress("./temp/" + filename, "./temp/hlae"); err != nil {
		return err
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在转移文件...")
	//转移
	if err := copy.Copy("./temp/hlae", a.cfg.HlaePath); err != nil {
		return err
	}
	_ = os.RemoveAll("./temp")

	//生成version文件
	a.cfg.HlaeVersion = version
	ver := strings.Replace(version, "v", "", -1)
	if err := tool.WriteFast(a.cfg.HlaePath + "/version", ver); err != nil {
		return err
	}

	//完成
	a.setProgress(100)
	a.setLog("HLAE更新完成")
	return nil
}

//更新FFmpeg
func (a *App) updateFFmpeg() error {
	//声明变量
	var srcVersion, cdnVersion, srcURL, cdnURL, srcFilename, cdnFilename, version, url, filename string

	//声明用来控制进度条的变量
	const count = 5
	i := 0
	a.setProgress(0)
	a.setLog("正在读取FFmpeg和CDN源API...")
	//尝试读取FFmpeg和CDN源API
	srcData, err1 := tool.GetHttpData(a.cfg.FFmpegAPI + "/release-version")
	cdnData, err2 := tool.GetHttpData(a.cfg.FFmpegCdnAPI)
	if err1 != nil && err2 != nil {
		return errors.New("FFmpeg官方和CDN的API均获取失败")
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在解析API...")
	//解析API信息，决定下载的文件
	if err1 == nil {
		//解析官方API
		srcVersion = srcData
		srcURL = a.cfg.FFmpegAPI + "/ffmpeg-release-essentials.7z"
		srcFilename = "ffmpeg-release-essentials.7z"
	}

	if err2 == nil {
		//解析CDN API
		var cdnInst api.ReleaseDelivr
		//注释下面一行->使用encoding/json库
		var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary //使用高性能json-iterator/go库
		err := jsonx.Unmarshal([]byte(cdnData), &cdnInst)      //第二个参数要地址传递
		if err != nil {
			return err
		}
		//获取版本号、下载地址和文件名
		cdnVersion = cdnInst.Version
		cdnURL = cdnInst.DownloadLink[0]
		_, cdnFilename = filepath.Split(cdnURL)
	}

	//决定下载的文件
	if srcVersion == "" && cdnVersion == "" {
		return errors.New("FFmpeg官方和CDN的API均获取或解析失败")
	} else if srcVersion == "" {
		a.noticeWarning("FFmpeg官方API解析失败，CDN源可能不是最新版本")
	} else if cdnVersion == "" {
		a.noticeWarning("CDN源解析失败，下载速度可能较慢")
	}
	if srcVersion == "" || (srcVersion != "" && srcVersion == cdnVersion) {
		//官方版本非空且和CDN源版本一致，则下载CDN源
		version = cdnVersion
		url = cdnURL
		filename = cdnFilename
	} else {
		//否则下载官方源
		version = srcVersion
		url = srcURL
		filename = srcFilename
	}

	//对比版本号，决定是否更新
	//a.noticeWarning(version + " ? " + a.cfg.FFmpegVersion)
	if version == a.cfg.FFmpegVersion {
		return nil
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在下载FFmpeg安装包...")
	//下载FFmpeg安装包
	if err := tool.DownloadFile(url, "./temp"); err != nil {
		a.noticeError("FFmpeg下载失败")
		return err
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在解压FFmpeg安装包...")
	//解压
	if err := tool.Decompress("./temp/" + filename, "./temp/ffmpeg"); err != nil {
		return err
	}

	i += 100
	a.setProgress(i / count)
	a.setLog("正在转移文件...")

	//转移
	if ok, err := tool.IsFileExisted("./temp/ffmpeg/bin/ffmpeg.exe"); err != nil || !ok {
		if ok, err := tool.IsFileExisted("./temp/ffmpeg/ffmpeg-release-essentials/bin/ffmpeg.exe"); err != nil || !ok {
			return err
		} else {
			if err := copy.Copy("./temp/ffmpeg/ffmpeg-release-essentials", a.cfg.HlaePath + "/ffmpeg"); err != nil {
				return err	//\temp\ffmpeg\ffmpeg-release-essentials\bin
			}
		}
	} else {
		if err := copy.Copy("./temp/ffmpeg", a.cfg.HlaePath + "/ffmpeg"); err != nil {
			return err
		}
	}
	_ = os.RemoveAll("./temp")

	//完成
	a.cfg.FFmpegVersion = version
	a.setProgress(100)
	a.setLog("FFmpeg更新完成")
	return nil
}

//启动HLAE
func (a *App) LaunchHLAE() bool {
	if a.cfg.HlaeState == true {
		_ = a.runtime.Browser.OpenFile(a.cfg.HlaePath + "/hlae.exe")
	}

	return true
}

//打开HLAE文件夹
func (a *App) OpenHlaeDirectory() error {
	if a.cfg.HlaeState == true {
		if err := a.runtime.Browser.OpenFile(a.cfg.HlaePath); err != nil {
			a.setLog(err.Error())
			log.Println(err)
			return err
		}
	}
	return nil
}
