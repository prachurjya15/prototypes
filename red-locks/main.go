package main

import (
	"sync"

	"github.com/prachurjya15/prototypes/red-locks/consumers"
	"github.com/redis/go-redis/v9"
)

const N_CONSUMERS = 10
const QUORUM = 3

func createAndGetRedisClients(clientAddrs []string) []*redis.Client {
	var redisClients []*redis.Client
	for _, addr := range clientAddrs {
		rc := redis.NewClient(&redis.Options{
			Addr: addr,
		})
		redisClients = append(redisClients, rc)
	}
	return redisClients
}

func main() {
	redisCltsAddr := []string{"localhost:6379",
		"localhost:6380",
		"localhost:6381",
		"localhost:6382",
		"localhost:6383"} // Load it from env

	redisClients := createAndGetRedisClients(redisCltsAddr)

	cs := make([]consumers.Consumer, 0)
	for i := range N_CONSUMERS {
		c := consumers.NewConsumer(i+1, redisClients, QUORUM)
		cs = append(cs, c)
	}

	var wg sync.WaitGroup
	for _, c := range cs {
		c := c
		wg.Go(func() { c.PerformTask() })
	}
	wg.Wait()
}
