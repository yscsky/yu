package main

import (
	"context"

	"github.com/yscsky/yu"
	"github.com/yscsky/yu/examples/server_grpc_usage/pb"
	"google.golang.org/grpc"
)

func main() {
	gs := yu.NewGrpcServer("GrpcUse", ":8080", setRegister, grpc.UnaryInterceptor(yu.GrpcRecovery))
	yu.Run(&yu.App{
		Na:    "ServerGrpcUsage",
		Start: func() bool { return true },
		Stop:  func() {},
		Svrs:  []yu.ServerInterface{gs},
	})
}

func setRegister(gs *yu.GrpcServer) {
	pb.RegisterHelloServiceServer(gs, &server{})
}

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	yu.Logf("%s SayHello", req.Name)
	if req.Name == "panic" {
		panic(req.Name)
	}
	return &pb.HelloResponse{Reply: "Hello " + req.Name}, nil
}
