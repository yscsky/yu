package yu

import (
	"os"
	"os/signal"
	"runtime/debug"
)

// AppInterface 应用实现接口
type AppInterface interface {
	Name() string
	Servers() []ServerInterface
	OnStart() bool
	OnStop()
}

// ServerInterface 服务实现接口
type ServerInterface interface {
	OnStart() bool
	OnStop()
	Info() string
}

// Run 运行服务
func Run(app AppInterface) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	Logf("app %s start", app.Name())
	if !app.OnStart() {
		Logf("app run stop, start return false")
		return
	}
	for _, server := range app.Servers() {
		// 启动的服务有些会阻塞，因此放goroutine中执行
		go func(svr ServerInterface) {
			if !svr.OnStart() {
				Errf("server start failed: %s", svr.Info())
			}
		}(server)
	}
	sig := make(chan os.Signal)
	defer close(sig)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	for _, server := range app.Servers() {
		server.OnStop()
	}
	app.OnStop()
	Logf("app %s stop", app.Name())
}

// App 内置App
type App struct {
	Na    string
	Start func() bool
	Stop  func()
	Svrs  []ServerInterface
}

// Name 实现 Name() string 接口
func (a *App) Name() string {
	return a.Na
}

// Servers 实现 Servers() []ServerInterface 接口
func (a *App) Servers() []ServerInterface {
	return a.Svrs
}

// OnStart 实现 OnStart() bool 接口
func (a *App) OnStart() bool {
	if a.Start == nil {
		//! 没有设置Start不能启动
		return false
	}
	return a.Start()
}

// OnStop 实现 OnStop() 接口
func (a *App) OnStop() {
	if a.Stop == nil {
		return
	}
	a.Stop()
}
