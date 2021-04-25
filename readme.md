# yu

Golang 开发辅助库

# 导入

```go
import "github.com/yscsky/yu"
```

## 代码自动生成工具

见 [使用说明](cmd/yugen/readme.md)

## 服务

`app.go`提供了一个服务启动的框架，有监听中断信号的功能。

可以自行实现`AppInterface`和`ServerInterface`接口，也可以用库中提供的。

参考`examples/app_base_usage`中代码。

`server_gin.go`,`server_grpc.go`是分别实现`ServerInterface`接口的服务。

可以参考`examples`下的`server_gin_usage`，`server_grpc_usage`中代码。

`middleware.go`提供的 gin 和 grpc 中间件。

`client_http.go`和`client_grpc.go`分别提供了 http 和 grpc 的客户端生成函数，可以参考`examples`下的`client_gin_usage`，`client_grpc_usage`中代码。

`prometheus.go`添加 prometheus 监控，提供了 http 和 grpc 的中间件。

## 数据库

`db_gorm.go` 使用 gorm 连接数据库，内置了 mysql 和 postgres 的连接方式，参考`examples/db_gorm_usage`中代码。

`db_sqlx.go` 使用 sqlx 连接数据库，内置了一个加了 statement 缓存的数据库操作结构，参考`examples/db_sqlx_usage`中代码。

`db_redis.go` 提供连接 redis 的函数，使用`github.com/go-redis/redis/v8`包。

## 其他

`queue.go` 使用 chan 实现多协程的队列处理，参考`examples/queue_usage`中代码。

`serial.go` 使用 chan 实现线程安全的序列生成器，参考`examples/serial_usage`中代码。

`tstamp.go` 实现一个时间戳类型，JSON 序列化为字符串，反序列化为 int64。

`utils.go` 提供了一系列便利的函数。
