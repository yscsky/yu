# 基于 yu 的代码自动生成器

## 使用参数说明

### -m string

设置 go mod init 的名称，一般为项目的 git 路径，必须设置。

### -b bool

是否创建完脚本后进行 go build，默认执行。

## 生成文件目录

```shell
.
├── Dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   └── app.go
│   ├── db
│   │   ├── db.go
│   │   └── sqls.go
│   └── model
│       ├── config.go
│       ├── const.go
│       └── model.go
├── main.go
└── readme.md
```
