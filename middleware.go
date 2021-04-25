package yu

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// NoCache 设置请求头为无缓存
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		c.Next()
	}
}

// BasicAuth 设置BasicAuth认证
func BasicAuth(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if u, p, ok := c.Request.BasicAuth(); ok && u == username && p == password {
			c.Next()
			return
		}
		c.Writer.Header().Set("WWW-Authenticate", `Basic realm="Authorization Required"`)
		http.Error(c.Writer, "unauthorized", http.StatusUnauthorized)
		log.Println("[ERR] - unauthorized")
		c.Abort()
	}
}

// LogControl 控制gin是否输出log
func LogControl(trace bool, skip []string) gin.HandlerFunc {
	logger := gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: skip})
	return func(c *gin.Context) {
		if !trace {
			c.Next()
			return
		}
		logger(c)
	}
}

// CORS 设置跨域
func CORS(acceptedOrigin, AcceptedMethods, AcceptedHeaders, MaxAge string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", acceptedOrigin)
		if AcceptedMethods != "" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", AcceptedMethods)
		}
		if AcceptedHeaders != "" {
			c.Writer.Header().Set("Access-Control-Allow-Headers", AcceptedHeaders)
		}
		if MaxAge != "" {
			c.Writer.Header().Set("Access-Control-Max-Age", MaxAge)
		}
		if c.Request.Method != "OPTIONS" {
			c.Next()
		}
	}
}

// GrpcRecovery grpc接口panic恢复中间件
func GrpcRecovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	return handler(ctx, req)
}
