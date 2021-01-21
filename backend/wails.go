package backend

import "fmt"

///// wails.go存放backend包对frontend细粒度交互的操作

//测试发送信息
func (a *App) SayHello() string {
	fmt.Println("Hello to Backend!")
	//发信息 事件.信号  Emit
	//a.runtime.Events.Emit("error", "Go发送的错误信息！", 6657)
	return "Hello to Frontend!"
}

//设置进度条 需要前端Mount
func (a *App) setProgressBar(percent int) {
	a.runtime.Events.Emit("SetProcessBar", percent)
}

