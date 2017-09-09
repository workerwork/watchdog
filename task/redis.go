package task

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	//RedisClient  *redis.Pool
	REDIS_SERVER = "127.0.0.1:6379"
	//REDIS_PASSWORD     = ""
	REDIS_DB int = 0
)

func newPool(server string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			/*if _, err := c.Do("AUTH", password); err != nil {
			    c.Close()
								return nil, err
											}*/
			c.Do("SELECT", db)
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func writeRedis(cmd string, args ...interface{}) {
	pool := newPool(REDIS_SERVER, REDIS_DB)
	rs := pool.Get()
	if _, err := rs.Do(cmd, args...); err != nil {
		fmt.Println(err)
	}
	defer rs.Close()
}

func readRedis(cmd string, args ...interface{}) (value []string) {
	pool := newPool(REDIS_SERVER, REDIS_DB)
	rs := pool.Get()
	value, err := redis.Strings(rs.Do(cmd, args...))
	if err != nil {
		fmt.Println(err)
	}
	defer rs.Close()
	return
}

func existsRedis(key string) (exists bool) {
	pool := newPool(REDIS_SERVER, REDIS_DB)
	rs := pool.Get()
	exists, err := redis.Bool(rs.Do("EXISTS", key))
	if err != nil {
		fmt.Println(err)
	}
	defer rs.Close()
	return
}
