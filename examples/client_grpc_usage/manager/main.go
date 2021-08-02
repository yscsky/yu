package main

import (
	"context"
	"time"

	"github.com/yscsky/yu"
	"github.com/yscsky/yu/examples/server_grpc_usage/pb"
	"google.golang.org/grpc"
)

type Client struct {
	pb.HelloServiceClient
}

func main() {
	manager := yu.NewGrpcConnManager()
	defer manager.CloseConns()

	client := &Client{}
	manager.AddConn("hello", ":8080", client)

	time.Sleep(time.Second)

	reply, err := client.say("God")
	if err != nil {
		yu.LogErr(err, "say")
		return
	}
	yu.Logf("reply: %s", reply)
}

func (c *Client) OnConnected(conn *grpc.ClientConn) {
	c.HelloServiceClient = pb.NewHelloServiceClient(conn)
	yu.Logf("NewHelloServiceClient ok")
}

func (c *Client) say(name string) (reply string, err error) {
	resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		return
	}
	reply = resp.Reply
	return
}
