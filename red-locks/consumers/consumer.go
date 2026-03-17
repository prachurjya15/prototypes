package consumers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Consumer struct {
	Consumer_id  string
	RedisClients []*redis.Client
	Quorum       int
	expiration   time.Duration
}

// lua fn : get key, if !key set and return true else if val != consumer_id return false else return true

const luaLockScript = `
local val = redis.call('GET', KEYS[1])
if not val then
	redis.call('SET', KEYS[1], ARGV[1], 'PX', ARGV[2])
	return 1
elseif val == ARGV[1] then
	redis.call('PEXPIRE', KEYS[1], ARGV[2])
	return 1
else return 0
end
`

const luaUnlockScript = `
local val = redis.call('GET', KEYS[1])
if not val then
	return
end
if val == ARGV[1] then
	redis.call('DEL', KEYS[1])
	return
end
`

func NewConsumer(index int, r []*redis.Client, quorum int) Consumer {
	c := Consumer{
		RedisClients: r,
		Quorum:       quorum,
		expiration:   11 * time.Second,
	}
	c.Consumer_id = fmt.Sprintf("Consumer-%d", index)
	return c
}

func (c *Consumer) PerformTask() {
	// 1.Take Lock
	c.lock()
	// 2. Perform Task
	c.doWork()
	// 3.Release Lock
	c.unLock()
}

func acquireLock(r *redis.Client, consumerId string, ttl time.Duration) bool {
	retCode, err := r.Eval(context.Background(),
		luaLockScript,
		[]string{"LOCK"},
		consumerId,
		ttl.Milliseconds(),
	).Int()
	if err != nil {
		log.Printf("Error in executing the lock lua script. Err: [%s]\n", err)
		return false
	}
	return retCode == 1
}

func releaseLock(r *redis.Client, consumerId string) {
	err := r.Eval(context.Background(),
		luaUnlockScript,
		[]string{"LOCK"},
		consumerId,
	).Err()
	if err != nil && err != redis.Nil {
		log.Printf("Error in executing the unlock lua script. Error: [%s]\n", err)
		return
	}
}

func (c *Consumer) lock() {
	// We will iterate over the redis clients.
	// Get value of Key: Lock and check if it is held by any consumers.
	// If its the curr consumer holding the lock then mark as acquired and proceed.
	// If its not the curr consumer then wait for 1 sec and retry
	// If no one is holding the key then try to put in Key: Lock and Val: <consumer_id>
	for {
		nAck := 0
		nAckServer := make([]*redis.Client, 0)
		for _, rc := range c.RedisClients {
			acquired := acquireLock(rc, c.Consumer_id, c.expiration)
			if acquired {
				nAck++
				nAckServer = append(nAckServer, rc)
			}
		}
		if nAck >= c.Quorum {
			break
		} else {
			// Release all the locks held
			for _, srv := range nAckServer {
				releaseLock(srv, c.Consumer_id) // TODO: Will need retry mechanisms for unreliable network or error in cmds
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (c *Consumer) unLock() {
	for _, rc := range c.RedisClients {
		releaseLock(rc, c.Consumer_id)
	}
}

func (c *Consumer) doWork() {
	log.Printf("##### %s got lock and started working ###### \n", c.Consumer_id)
	time.Sleep(10 * time.Second)
	log.Printf("##### %s done working ###### \n", c.Consumer_id)
}
