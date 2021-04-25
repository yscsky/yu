package main

import (
	"sync"

	"github.com/yscsky/yu"
)

var wg *sync.WaitGroup

func main() {
	wg = new(sync.WaitGroup)
	wg.Add(1)
	// 一种方式是实现AppInterface接口
	go yu.Run(&App{name: "AppBaseSelf"})
	wg.Add(1)
	// 另一种方式是使用帮助库提供的内置App
	go yu.Run(&yu.App{
		Na: "AppBaseDefault",
		Start: func() bool {
			yu.Logf("设置Start，执行初始化操作")
			return true
		},
		Stop: func() {
			yu.Logf("设置Stop，执行停止操作")
			wg.Done()
		},
		Svrs: []yu.ServerInterface{},
	})
	wg.Wait()
}

type App struct {
	name string
}

func (a *App) Name() string {
	return a.name
}

func (a *App) OnStart() bool {
	yu.Logf("这里执行数据库连接，定时任务设置，客户端初始化，全局变量初始化，队列初始化，各种服务初始化等操作")
	return true
}

func (a *App) OnStop() {
	yu.Logf("这里执行关闭数据库，停止队列等操作")
	wg.Done()
}

func (a *App) Servers() []yu.ServerInterface {
	return []yu.ServerInterface{}
}
