package yu

import (
	"net"

	"google.golang.org/grpc"
)

// GrpcServer 封装grpc server
type GrpcServer struct {
	Name string
	Addr string
	*grpc.Server
	Register func(*GrpcServer)
}

// NewGrpcServer 创建GrpcServer
func NewGrpcServer(name, addr string, reg func(*GrpcServer), opt ...grpc.ServerOption) *GrpcServer {
	return &GrpcServer{
		Name:     name,
		Addr:     addr,
		Server:   grpc.NewServer(opt...),
		Register: reg,
	}
}

// OnStart 实现ServerInterface接口
func (gs *GrpcServer) OnStart() bool {
	if gs.Register == nil {
		Errf("GrpcServer Register is nil")
		return false
	}
	gs.Register(gs)
	lis, err := net.Listen("tcp", gs.Addr)
	if err != nil {
		LogErr(err, "net.Listen: "+gs.Addr)
		return false
	}
	Logf("%s grpc server start at %s", gs.Info(), gs.Addr)
	if err = gs.Serve(lis); err != nil {
		LogErr(err, "Serve")
	}
	return true
}

// OnStop 实现ServerInterface接口
func (gs *GrpcServer) OnStop() {
	gs.GracefulStop()
}

// Info 实现ServerInterface接口
func (gs *GrpcServer) Info() string {
	return gs.Name
}
