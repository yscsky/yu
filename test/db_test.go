package test

import (
	"context"
	"testing"

	"github.com/yscsky/yu"
)

func TestNewRedisClient(t *testing.T) {
	ctx := context.Background()
	cli := yu.NewRedisClient("127.0.0.1:6379", "", 0)
	pong, err := cli.Ping(ctx).Result()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pong)

	clients, err := yu.NewRedisClients("127.0.0.1:6379", "")
	if err != nil {
		t.Error(err)
		return
	}
	for db, cli := range clients {
		pong, err := cli.Ping(ctx).Result()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(db, pong)
	}
}
