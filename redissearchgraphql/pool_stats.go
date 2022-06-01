package redissearchgraphql

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func PoolStats(pool *redis.Pool) error {
	for {
		time.Sleep(time.Second * 10)
		promPoolActive.Set(float64(pool.Stats().ActiveCount))
		promPoolIdle.Set(float64(pool.Stats().IdleCount))
		promPoolWait.Set(float64(pool.Stats().WaitCount))
	}
}
