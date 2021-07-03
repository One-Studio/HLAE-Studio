package backend

import (
	"HLAE-Studio/backend/api"
	"HLAE-Studio/backend/config"
	"errors"
	"fmt"
	tool "github.com/One-Studio/ptools/pkg"
	"github.com/otiai10/copy"
	"github.com/wailsapp/wails"
	"log"
	"math/rand"
	"os"
	"os/user"
	"strings"
	"time"
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

	fmt.Println("wails测试")
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
	a.setProgress(0)
	a.setLog("正在初始化")
	if err = a.updateApp(); err != nil {
		a.noticeError(err.Error())
	}

	//检查是否初始化，否则通知前端选择安装方式和位置
	if a.cfg.Init == false {
		a.doSelectOption()
		return
	}

	//检查HLAE和FFmpeg安装情况
	a.setProgress(10)
	a.setLog("正在检查安装情况")
	if err = a.checkState(); err != nil {
		a.runtime.Events.Emit("SetLog", err)
		log.Println(err)
		return
	}

	//安装或者更新HLAE
	a.setProgress(20)
	a.setLog("正在下载/更新HLAE")
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
	a.setProgress(60)
	a.setLog("正在下载/更新FFmpeg")
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
	_ = os.RemoveAll("./temp")
	return
}

//检查HLAE和FFmpeg安装情况
func (a *App) checkState() error {
	//如果已经初始化，再核实一下hlae和ffmpeg程序位置
	if a.cfg.Init == true {
		//检查hlae
		ok1 := tool.IsFileExisted(a.cfg.HlaePath + "/hlae.exe")
		ok2 := tool.IsFileExisted(a.cfg.HlaePath + "/changelog.xml")
		if ok1 && ok2 {
			a.cfg.HlaeState = true
			//解析修正本地hlae版本号
			if ver, err := a.getLocalHlaeVersion(); err != nil {
				return err
			} else if ver != "" {
				a.cfg.HlaeVersion = ver
			}

			//检查ffmpeg
			if ffmpegOK := tool.IsFileExisted(a.cfg.HlaePath + "/ffmpeg/bin/ffmpeg.exe"); ffmpegOK {
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

//获取本地HLAE版本号，如果未安装则返回空版本号
func (a *App) getLocalHlaeVersion() (string, error) {
	//识别已安装hlae，则检查更新
	if ok := tool.IsFileExisted(a.cfg.HlaePath + "/hlae.exe"); !ok {
		//本地文件不存在，返回空版本号
		return "", nil
	} else {
		//获取本地版本号
		if changelog, err := tool.ReadAll(a.cfg.HlaePath + "/changelog.xml"); err != nil {
			return "", err
		} else if version, err := api.ParseChangelog(changelog); err != nil {
			return "", err
		} else {
			return version, nil
		}
	}
}

//获取最新better-github-api某工具的版本号
func (a *App) getVersion(api string) (version string, err error) {
	for count := 0;count < 5; count++ {
		version, err = tool.GetHttpData(api + "/version")
		if err != nil {
			time.Sleep(2 * time.Second)
		} else {
			continue
		}
	}

	return
}

//生成随机字符串
func (a *App) randomString(n int) string  {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%X", randBytes)
}

//安装HLAE
func (a *App) installHLAE() error  {
	//根据关联安装与否更新路径
	if a.cfg.Standalone == true {
		a.noticeWarning("请选择安装位置，或者已有hlae.exe的文件夹")
		//time.Sleep(2 * time.Second)
		path := a.SelectDirectory()
		if path == "" {
			a.noticeWarning("已取消安装")
			return nil
		}

		//更新hlae路径
		a.cfg.HlaePath = tool.FormatPath(path)
	} else {
		//检查CSGO Demos Manager是否安装 "%HOMEDIR%/AppData/Local/AkiVer/"
		usr, err := user.Current()
		if err != nil {
			return err
		}
		if ok := tool.IsFileExisted(usr.HomeDir + "/AppData/Local/AkiVer/");!ok {
			//Manager未安装  提示安装且弹出下载的网页
			if err := a.runtime.Browser.OpenURL("https://github.com/akiver/CSGO-Demos-Manager/releases/latest"); err != nil {
				return err
			}
			if err := a.runtime.Browser.OpenURL("https://api.upup.cool/get/csdm"); err != nil {
				return err
			}

			return errors.New("CSGO Demos Manager未安装，已弹出下载地址")
		}

		//更新hlae路径
		a.cfg.HlaePath = tool.FormatPath(usr.HomeDir + "/AppData/Local/AkiVer/hlae")
	}

	//确定位置后调用更新HLAE的方法
	err := a.updateHLAE()
	return err
}

//更新HLAE
func (a *App) updateHLAE() error {
	//获取版本
	a.setProgress(30)
	a.setLog("正在获取HLAE最新版本")
	var localVersion, latestVersion string
	var err error
	localVersion = a.cfg.HlaeVersion
	if latestVersion, err = a.getVersion(a.cfg.HlaeAPI); err != nil {
		return err
	}
	if localVersion == latestVersion {
		return nil
	}

	//下载HLAE安装包
	a.setProgress(40)
	a.setLog("正在下载HLAE安装包")
	var filepath = "./temp/" + a.randomString(12)
	if err = tool.DownloadFile(filepath, a.cfg.HlaeAPI); err != nil {
		a.noticeError("HLAE下载失败")
		return err
	}

	//获取文件路径+名
	var file = tool.GetFilePathFromDir(filepath, "hlae")

	//解压
	a.setProgress(50)
	a.setLog("正在解压HLAE")
	if err = tool.Decompress(file, "./temp/hlae"); err != nil {
		return err
	}

	//转移
	if err = copy.Copy("./temp/hlae", a.cfg.HlaePath); err != nil {
		return err
	}
	_ = os.RemoveAll("./temp/hlae")
	_ = os.RemoveAll(filepath)
	a.cfg.HlaeVersion = latestVersion

	//生成version文件
	if !a.cfg.Standalone {
		ver := strings.Replace(a.cfg.HlaeVersion, "v", "", -1)
		if err = tool.WriteFast(a.cfg.HlaePath + "/version", ver); err != nil {
			return err
		}
	}

	return nil
}

//安装FFmpeg
func (a *App) installFFmpeg() error {
	return a.updateFFmpeg()
}

//更新FFmpeg
func (a *App) updateFFmpeg() error {
	//获取版本
	a.setProgress(70)
	a.setLog("正在获取FFmpeg最新版本")
	var localVersion, latestVersion string
	var err error
	localVersion = a.cfg.FFmpegVersion
	if latestVersion, err = a.getVersion(a.cfg.FFmpegAPI); err != nil {
		return err
	}
	if localVersion == latestVersion {
		return nil
	}

	//下载FFmpeg安装包
	a.setProgress(80)
	a.setLog("正在下载FFmpeg")
	var filepath = "./temp/" + a.randomString(12)
	if err = tool.DownloadFile(filepath, a.cfg.FFmpegAPI); err != nil {
		a.noticeError("FFmpeg下载失败")
		return err
	}

	//获取文件路径+名
	var file = tool.GetFilePathFromDir(filepath, "ffmpeg")

	//解压
	a.setProgress(90)
	a.setLog("正在解压FFmpeg")
	if err = tool.Decompress(file, "./temp/ffmpeg"); err != nil {
		return err
	}

	//转移
	file = tool.GetFilePathFromDir("./temp/ffmpeg", "ffmpeg.exe")
	if file == "" {
		return errors.New("找不到ffmpeg解压后的程序文件")
	}
	if err = copy.Copy(file, a.cfg.HlaePath + "/ffmpeg/bin/ffmpeg.exe"); err != nil {
		return err
	}

	_ = os.RemoveAll("./temp/ffmpeg")
	_ = os.RemoveAll(filepath)
	a.cfg.FFmpegVersion = latestVersion
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
