package backend

import (
	"HLAE-Studio/backend/config"
	"fmt"
	"github.com/wailsapp/wails"
)

///// app.go 存放backend包与frontend交互的大部分操作

//App设置
type App struct {
	runtime *wails.Runtime //初始化Runtime需要
	cfg     config.CFG
	//zoom
}

//Wails初始化
func (a *App) WailsInit(runtime *wails.Runtime) error {
	a.runtime = runtime
	//初始化后：
	fmt.Println("Wails初始化")
	config.ReadConfig(a.cfg)

	return nil
}

//Wails结束前
func (a *App) WailsShutdown() {
	//结束前：
	fmt.Println("Wails结束")
	config.SaveConfig(a.cfg)

	return
}

func (a *App) SelectSrcPath() (string, error) {

	return "", nil
}

func (a *App) SelectDstPath() (string, error) {

	return "", nil
}

func (a *App) ParseDragFiles() (string, error) {

	return "", nil
}

func (a *App) StartEncoding() (string, error) {

	return "", nil
}

func (a *App) PauseEncoding() (string, error) {

	return "", nil
}

func (a *App) CheckUpdate() error {

	return nil
}

func (a *App) BrowseLog() (string, error) {

	return "", nil
}

func (a *App) OpenProgramDir() (string, error) {

	return "", nil
}
