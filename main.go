package main

import (
	"HLAE-Studio/backend"
	_ "embed"
	"github.com/wailsapp/wails"
)

//绑定js和css

//go:embed frontend/dist/app.js
var js string

//go:embed frontend/dist/app.css
var css string

func main() {
	//创建app 设定宽高标题颜色和是否锁定长宽
	app := wails.CreateApp(&wails.AppConfig{
		Width:  350,
		Height: 525, //windows标题栏40px 总宽度566.34px
		Title:  "HLAE Studio",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
		//DisableInspector: false,
		//Resizable: true,
	})

	//绑定后端&运行
	app.Bind(&backend.App{})
	_ = app.Run()
}
