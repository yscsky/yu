package yu

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// NewGrpcConn 创建grpc.ClientConn，with insecure and block
func NewGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
}

// TraceUnaryInt 计算客户端请求耗时
func TraceUnaryInt(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("[INFO] - %s exec %s", method, time.Since(start))
	return err
}
