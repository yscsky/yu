package main

import (
	"context"
	"flag"

	"github.com/yscsky/yu"
	"github.com/yscsky/yu/examples/server_grpc_usage/pb"
)

func main() {
	msg := flag.String("msg", "", "send message")
	flag.Parse()
	conn, err := yu.NewGrpcConn(":8080")
	if err != nil {
		yu.LogErr(err, "NewGrpcConn")
		return
	}
	defer conn.Close()
	client := pb.NewHelloServiceClient(conn)
	if *msg != "" {
		call(client, *msg)
		return
	}
	for i := 0; i < 100; i++ {
		call(client, "client")
	}
}

func call(client pb.HelloServiceClient, msg string) {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: msg})
	if err != nil {
		yu.LogErr(err, "SayHello")
		return
	}
	yu.Logf(resp.Reply)
}
