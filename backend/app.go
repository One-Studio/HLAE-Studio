package backend

import (
	"HLAE-Studio/backend/config"
	"HLAE-Studio/backend/tool"
	"fmt"
	"github.com/wailsapp/wails"
	"log"
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

	//设置前端
	a.setAppVersion(a.cfg.AppVersion)
	a.setVersionCode(a.cfg.VersionCode)

	return nil
}

//Wails结束前
func (a *App) WailsShutdown() {
	//结束前：
	fmt.Println("Wails结束")
	err := config.SaveConfig(a.cfg,"./config.json")
	if err != nil {
		log.Println(err)
	}
}

//检查更新
func (a *App) CheckUpdate() error {
	//检查HLAE和FFmpeg安装情况
	if err := a.checkState(); err != nil {
		a.runtime.Events.Emit("SetLog", err)
		log.Println(err)
		return err
	}
	//TODO



	return nil
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
		} else {
			a.cfg.HlaeState = false
		}
		//检查ffmpeg
		if ok, err := tool.IsFileExisted(a.cfg.HlaePath + "/ffmpeg/bin/ffmpeg.exe"); err != nil {
			return err
		} else if ok == true {
			a.cfg.FFmpegState = true
		} else {
			a.cfg.FFmpegState = false
		}
	} else {
		//否则直接给定false开始安装hlae+ffmpeg
		a.cfg.HlaeState = false
		a.cfg.FFmpegState = false
		if err := a.selectOption(); err != nil {
			return err
		}
		if err := a.installHLAE(); err != nil {
			return err
		}
		if err := a.installFFmpeg(); err != nil {
			return err
		}
	}

	return nil
}

//TODO 选择安装hlae的方式
func (a *App) selectOption() error {

	return nil
}

//安装HLAE TODO
func (a *App) installHLAE() error {
	//声明用来控制进度条的变量
	const count = 5
	i := 0
	a.setProgress(i/count)
	a.setLog("正在读取HLAE和CDN源API...")
	i+=100
	//尝试读取HLAE和CDN源API


	a.setProgress(i/count)
	a.setLog("正在解析API...")
	i+=100
	//解析API信息，决定下载的文件


	a.setProgress(i/count)
	a.setLog("正在下载HLAE安装包...")
	i+=100
	//下载HLAE安装包


	a.setProgress(i/count)
	a.setLog("正在解压HLAE安装包...")
	i+=100
	//解压


	a.setProgress(i/count)
	a.setLog("正在转移文件...")
	i+=100
	//转移+生成version文件


	a.setProgress(i/count)
	a.setLog("HLAE安装完成")
	i+=100

	return nil
}

//安装FFmpeg
func (a *App) installFFmpeg() error {

	return nil
}

//启动HLAE
func (a *App) LaunchHLAE() bool {
	if a.cfg.HlaeState == true {
		out, err := tool.Cmd("explorer " + a.cfg.HlaePath + "/hlae.exe")
		if err != nil {
			a.setLog(out)
			log.Println(out, err)
			return false
		}
		return true
	}
	return false
}

//打开HLAE文件夹
func (a *App) OpenHlaeDirectory() error {

	return nil
}
