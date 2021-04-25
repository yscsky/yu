package yu

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient 根据db编号创建redis.Client
func NewRedisClient(addr, pass string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})
}

// NewRedisClients 创建redis.Client map以db编号为key
func NewRedisClients(addr, pass string) (clients map[int]*redis.Client, err error) {
	clients = make(map[int]*redis.Client)
	clients[0] = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})
	res, err := clients[0].ConfigGet(context.Background(), "databases").Result()
	if err != nil {
		return
	}
	if len(res) < 2 {
		err = errors.New("couldn't get databases")
		return
	}
	val, _ := res[1].(string)
	num, _ := strconv.Atoi(val)
	for i := 1; i < num; i++ {
		clients[i] = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass,
			DB:       i,
		})
	}
	return
}
