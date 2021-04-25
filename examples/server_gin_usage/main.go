package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yscsky/yu"
)

func main() {
	// 创建gin server
	gs := yu.NewGinServer("GinUse", ":8080", gin.DebugMode)
	// 设置健康检查
	gs.Health()
	gs.Promethous("admin", "admin")
	group := gs.Group("api", yu.NoCache(), yu.PromMetrics(), yu.LogControl(true, []string{}))
	group.GET("getjson", getJson)
	group.GET("getjsonauth", yu.BasicAuth("admin", "123456"), getJsonAuth)
	group.POST("postjson", postJson)
	group.POST("postjsonauth", yu.BasicAuth("admin", "123456"), postJsonAuth)
	group.POST("postform", postForm)
	group.POST("postformauth", yu.BasicAuth("admin", "123456"), postFormAuth)
	// 运行服务
	yu.Run(&yu.App{
		Na:    "ServerGinUsage",
		Start: func() bool { return true },
		Stop:  func() {},
		Svrs:  []yu.ServerInterface{gs},
	})
}

func getJson(c *gin.Context) {
	yu.JsonOK(c, "get json")
}

func getJsonAuth(c *gin.Context) {
	yu.JsonOK(c, "get json with auth")
}

func postJson(c *gin.Context) {
	var req yu.Resp
	if err := c.BindJSON(&req); err != nil {
		yu.JsonErr(c, err)
		return
	}
	yu.JsonOK(c, req.Data)
}

func postJsonAuth(c *gin.Context) {
	var req yu.Resp
	if err := c.BindJSON(&req); err != nil {
		yu.JsonErr(c, err)
		return
	}
	yu.JsonOK(c, req.Data)
}

func postForm(c *gin.Context) {
	data, ok := c.GetPostForm("data")
	if !ok {
		yu.JsonMsg(c, "data not found")
		return
	}
	yu.JsonOK(c, data)
}

func postFormAuth(c *gin.Context) {
	data, ok := c.GetPostForm("data")
	if !ok {
		yu.JsonMsg(c, "data not found")
		return
	}
	yu.JsonOK(c, data)
}
