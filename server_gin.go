package yu

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GinServer 包含gin的http server
type GinServer struct {
	Name   string
	Server *http.Server
	Engine *gin.Engine
}

// NewGinServer 创建GinServer
func NewGinServer(name, addr, mod string) *GinServer {
	gin.SetMode(mod)
	gs := &GinServer{
		Name:   name,
		Server: &http.Server{Addr: addr},
		Engine: gin.New(),
	}
	gs.Engine.Use(gin.Recovery(), gin.ErrorLogger())
	gs.Server.Handler = gs.Engine
	return gs
}

// OnStart 实现ServerInterface接口
func (s *GinServer) OnStart() bool {
	Logf("%s gin server start at %s", s.Info(), s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		LogErr(err, "ListenAndServe")
	}
	return true
}

// OnStop 实现ServerInterface接口
func (s *GinServer) OnStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		LogErr(err, "Shutdown")
	}
}

// Info 实现ServerInterface接口
func (s *GinServer) Info() string {
	return s.Name
}

// Health 健康检查接口，需要的话设置
func (s *GinServer) Health() {
	s.Engine.GET("health", func(c *gin.Context) { JsonOK(c, "yes") })
}

// Promethous 启动Promethous监听
func (s *GinServer) Promethous(name, pass string) {
	InitPrometheus(s.Name)
	s.Engine.GET("metrics", BasicAuth(name, pass), PromethousHandler())
}

// Group 创建gin.RouterGroup
func (s *GinServer) Group(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return s.Engine.Group(path, handlers...)
}

// JSON 发送JSON
func JSON(c *gin.Context, data interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err := json.NewEncoder(c.Writer).Encode(data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GzipJSON gzip压缩发送JSON
func GzipJSON(c *gin.Context, data interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")
	c.Writer.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(c.Writer)
	defer gz.Close()
	if err := json.NewEncoder(gz).Encode(data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

// JsonOK 反馈OK
func JsonOK(c *gin.Context, data interface{}) {
	JSON(c, RespOK(data))
}

// JsonErr 反馈Err
func JsonErr(c *gin.Context, err error) {
	JSON(c, RespErr(err))
}

// JsonMsg 反馈发生错误信息
func JsonMsg(c *gin.Context, msg string) {
	JSON(c, RespMsg(CodeErr, msg))
}
