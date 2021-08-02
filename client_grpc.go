package yu

import (
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// GrpcConnecter grpc连接接口
type GrpcConnecter interface {
	OnConnected(*grpc.ClientConn)
}

// GrpcConn grpc连接执行
type GrpcConn struct {
	addr       string
	clientConn *grpc.ClientConn
	connecter  GrpcConnecter
	opts       []grpc.DialOption
}

// GrpcConnManager grpc连接管理器
type GrpcConnManager struct {
	connMap map[string]*GrpcConn
	mux     *sync.Mutex
}

// NewGrpcConnManager 创建grpc连接管理器
func NewGrpcConnManager() *GrpcConnManager {
	return &GrpcConnManager{
		connMap: make(map[string]*GrpcConn),
		mux:     new(sync.Mutex),
	}
}

// NewGrpcConn 创建grpc.ClientConn，with insecure and block
func NewGrpcConn(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, append(opts, grpc.WithInsecure(), grpc.WithBlock())...)
}

// WithTrace 返回TraceUnaryInt的DialOption
func WithTrace() grpc.DialOption {
	return grpc.WithUnaryInterceptor(TraceUnaryInt)
}

// TraceUnaryInt 计算客户端请求耗时
func TraceUnaryInt(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("[INFO] - %s exec %s", method, time.Since(start))
	return err
}

// AddConn 添加grpc连接
func (gm *GrpcConnManager) AddConn(name, addr string, connecter GrpcConnecter, opts ...grpc.DialOption) {
	gc := &GrpcConn{
		addr:      addr,
		connecter: connecter,
		opts:      opts,
	}
	go gc.connect()
	gm.mux.Lock()
	defer gm.mux.Unlock()
	gm.connMap[name] = gc
}

// connect 发起grpc连接
func (gc *GrpcConn) connect() {
	gc.disconnect()
	conn, err := NewGrpcConn(gc.addr, gc.opts...)
	if err != nil {
		LogErr(err, "NewGrpcConn "+gc.addr)
		return
	}
	gc.clientConn = conn
	log.Printf("[INFO] - %s grpc connected", gc.addr)
	if gc.connecter != nil {
		gc.connecter.OnConnected(gc.clientConn)
	}
}

// disconnect 关闭grpc连接
func (gc *GrpcConn) disconnect() {
	if gc.clientConn != nil {
		gc.clientConn.Close()
	}
	gc.clientConn = nil
}

// CloseConns 关闭所有连接
func (gm *GrpcConnManager) CloseConns() {
	gm.mux.Lock()
	defer gm.mux.Unlock()
	for _, gc := range gm.connMap {
		gc.disconnect()
	}
}

// Reconnect 重新grpc连接
func (gm *GrpcConnManager) Reconnect(name string) {
	gm.mux.Lock()
	gc := gm.connMap[name]
	gm.mux.Unlock()
	gc.connect()
}
