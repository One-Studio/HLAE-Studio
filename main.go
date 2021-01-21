package main

import (
  "HLAE-Studio/backend"
  "github.com/leaanthony/mewn"
  "github.com/wailsapp/wails"
)

func main() {
  //绑定js和css
  js := mewn.String("./frontend/dist/app.js")
  css := mewn.String("./frontend/dist/app.css")

  //创建app 设定宽高标题颜色和是否锁定长宽
  app := wails.CreateApp(&wails.AppConfig{
    Width:  400,
    Height: 650,
    Title:  "HLAE-Studio",
    JS:     js,
    CSS:    css,
    Colour: "#131313",
    //Resizable: true,
  })

  //绑定后端&运行
  app.Bind(&backend.App{})
  app.Run()
}
