package redis

import (
	"context"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/gomodule/redigo/redis"
)

// ErrNil indicates that a reply value is nil.
var ErrNil = errors.New("(nil)")

// Redis describes connection to Redis server
type Redis struct {
	pool *redis.Pool
}

// New returns the initialized Redis object
func New(ctx context.Context, redisServer string) (*Redis, error) {
	log.Info("Redis init: host=", redisServer)
	c := new(Redis)
	c.initNewPool(redisServer)
	if err := c.Ping(); err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		err := c.pool.Close()
		if err != nil {
			log.Error("close redis connection error:", err.Error())
			return
		}
		log.Info("close redis connection")
	}()

	return c, nil
}

func (c *Redis) initNewPool(addr string) {
	c.pool = &redis.Pool{
		//MaxIdle:     5,
		//IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
		//TestOnBorrow: func(c redis.Conn, t time.Time) error {
		//	_, err := c.Do("PING")
		//	return err
		//},
		//MaxActive: 5,
		//Wait:      true,
	}
}

// Ping checks if connection exists
func (c *Redis) Ping() error {
	r := c.pool.Get()
	defer r.Close()
	_, err := r.Do("PING")
	return err
}
