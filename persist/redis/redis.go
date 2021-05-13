package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

// ConfigDef RedisConfig配置文件
type ConfigDef struct {
	Addr     string `json:"addr" yaml:"addr" mapstructure:"addr"`
	User     string `json:"user" yaml :"user" mapstructure:"user"`
	Password string `json:"password" yaml:"password" mapstructure:"password"`
	Db       int    `json:"db" yaml:"db" mapstructure:"db"`
	// Maximum number of idle connections in the pool.
	MaxIdle int `json:"maxIdle" yaml:"maxIdle" mapstructure:"maxIdle"`

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int `json:"maxActive" yaml:"maxActive" mapstructure:"maxActive"`

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout int64 `json:"idleTimeout" yaml:"idleTimeout" mapstructure:"idleTimeout"`

	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool `json:"wait" yaml:"wait" mapstructure:"wait"`

	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime int64 `json:"maxConnLifetime" yaml:"maxConnLifetime" mapstructure:"maxConnLifetime"`

	ReadTimeout    int64 `json:"readTimeout" yaml:"readTimeout" mapstructure:"readTimeout"`
	ConnectTimeout int64 `json:"connectTimeout" yaml:"connectTimeout" mapstructure:"connectTimeout"`
}

// NewRedisPool 通过配置文件初始化redis
func NewRedisPool(redisConfig *ConfigDef) (redisPool *Pool, err error) {
	rd := &redis.Pool{
		MaxIdle:     redisConfig.MaxIdle,
		MaxActive:   redisConfig.MaxActive,
		IdleTimeout: time.Duration(redisConfig.IdleTimeout * int64(time.Millisecond)),
		Wait:        redisConfig.Wait,
		//MaxConnLifetime: time.Duration(redisConfig.MaxConnLifetime*int64(time.Millisecond)),
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", redisConfig.Addr,
				redis.DialReadTimeout(time.Duration(redisConfig.ReadTimeout*int64(time.Millisecond))),
				redis.DialConnectTimeout(time.Duration(redisConfig.ConnectTimeout*int64(time.Millisecond))),
			)
			if err != nil {
				panic(fmt.Sprintf("redis.Dial err: %v, req: %v", err, redisConfig.Addr))
			}

			if auth := redisConfig.Password; auth != "" {
				if _, err = conn.Do("AUTH", auth); err != nil {
					panic(fmt.Sprintf("redis AUTH err: %v, req: %v,%v", err, redisConfig.Addr, auth))
				}
			}
			if db := redisConfig.Db; db > 0 {
				if _, err = conn.Do("SELECT", db); err != nil {
					panic(fmt.Sprintf("redis SELECT err: %v, req: %v,%v", err, redisConfig.Addr, db))
				}
			}
			return
		},
	}
	redisPool = &Pool{rd}
	err = checkRedis(rd)
	return
}

func checkRedis(redisPool *redis.Pool) error {
	var err error
	for i := 0; i < 3; i++ {
		err = checkRedisOnce(redisPool)
		if err == nil {
			return err
		}
	}
	return err
}

func checkRedisOnce(redisPool *redis.Pool) error {
	con := redisPool.Get()
	defer con.Close()
	_, err := con.Do("SETEX", "test", 1, 1)
	return err
}
