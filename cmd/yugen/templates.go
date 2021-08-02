package main

const gitignore = `.DS_Store
.idea/
data/
log/
vendor/
*.code-workspace
*.csv
*.db
*.log
{{.}}
`

const mainTemp = `package main

import (
	"{{.Mod}}/internal/app"
	"{{.Mod}}/internal/config"
	"github.com/yscsky/yu"
)

func init() {
	yu.CreateFolder("configs")
}

func main() {
	if err := config.NewConfig("configs/app.toml"); err != nil {
		yu.LogErr(err, "NewConfig")
		return
	}
	yu.Run(app.NewApp())
}
`

const appgo = `package app

import (
	"{{.Mod}}/internal/config"
	"{{.Mod}}/internal/db"
	"github.com/yscsky/yu"
)

var (
	ginSvr  *yu.GinServer
	grpcSvr *yu.GrpcServer
)

// App 实现AppInterface
type App struct{}

// NewApp 创建App
func NewApp() *App {
	return &App{}
}

// Name 实现AppInterface接口
func (a *App) Name() string {
	return "{{.Pn}}"
}

// Servers 实现AppInterface接口
func (a *App) Servers() []yu.ServerInterface {
	return []yu.ServerInterface{ginSvr, grpcSvr}
}

// OnStart 实现AppInterface接口
func (a *App) OnStart() bool {
	// 初始化数据库
	db.InitDB()

	// 设置gin server
	ginSvr = yu.NewGinServer(a.Name(), config.Cfg.HttpPort, config.Cfg.GinMod)
	ginSvr.Health()
	ginSvr.Promethous("admin", "admin")
	// group := ginSvr.Group("", yu.NoCache(), yu.PromeMetrics(), yu.LogControl(cfg.Trace, []string{}))

	// 设置grpc server
	grpcSvr = yu.NewGrpcServer(a.Name(), config.Cfg.GrpcPort, func(gs *yu.GrpcServer) {
	}, yu.PromeUnaryInterceptor())
	return true
}

// OnStop 实现AppInterface接口
func (a *App) OnStop() {
	// 关闭数据库
	db.CloseDB()
}
`

const dbgo = `package db

import (
	"{{.Mod}}/internal/config"
	"github.com/yscsky/yu"
)

var gdb *yu.GormDB

// InitDB 初始化数据库
func InitDB() {
	dsn := config.Cfg.MariaDSN
	gdb = yu.MustOpenMySQL(dsn)
	yu.Logf("mariadb open %s", dsn.MySQL())
	for _, table := range tables {
		if _, err := gdb.ExecSQL(table); err != nil {
			yu.LogErr(err, "create table")
			return
		}
	}
}

// CloseDB 关闭数据库
func CloseDB() {
	if gdb != nil {
		gdb.CloseDB()
	}
}
`
const sqlsgo = `package db

var tables = []string{}
`

const configgo = `package config

import (
	"github.com/yscsky/yu"
	"github.com/gin-gonic/gin"
)

// Config 项目配置
type Config struct {
	GrpcPort string
	HttpPort string
	GinMod   string
	Trace    bool
	MariaDSN yu.DSN
}

// Cfg 全局配置
var Cfg *Config

// NewConfig 创建项目配置
func NewConfig(path string) (err error) {
	Cfg = &Config{}
	err = yu.LoadOrSaveToml(path, Cfg, func() interface{} {
		Cfg = &Config{
			GrpcPort: ":8181",
			HttpPort: ":8080",
			GinMod:   gin.DebugMode,
			Trace:    true,
			MariaDSN: yu.DSN{
				Username: "root",
				Password: "123456",
				URL:      "127.0.0.1",
				Port:     "3306",
				DBName:   "dbname",
				SkipTran: false,
				PreStmt:  true,
				LogLevel: 0,
			},
		}
		return Cfg
	})
	return
}
`
const constgo = `package ml
`
const modelgo = `package ml
`

const dockerfile = `FROM registry.cn-hangzhou.aliyuncs.com/shortlog/go-alpine:1.16.5 as build

COPY . /{{.}}/

WORKDIR /{{.}}

ENV GOPROXY="https://goproxy.cn,direct" GOSUMDB="off"

RUN go mod tidy && go build

FROM registry.cn-hangzhou.aliyuncs.com/shortlog/alpcn:3.13

WORKDIR /app

COPY --from=build /{{.}}/{{.}} .

EXPOSE 8080 8181

CMD ["./{{.}}"]`

const readme = `# XX 项目

## 项目功能

## TODO

## 项目结构

## 数据库

## 接口说明
`
