package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/vmihailenco/msgpack"
)

// Pool redisPool
type Pool struct {
	*redis.Pool
}

// GetString get apiCache from redis.
func (rp *Pool) GetString(key string) (string, error) {
	conn := rp.Get()
	defer conn.Close()
	v, err := conn.Do("GET", key)
	if err != nil {
		return "", err
	}
	if v == nil {
		return "", nil
	}
	return redis.String(v, err)
}

// GetPrimitiveInterface Get redis from db.
func (rp *Pool) GetPrimitiveInterface(key string) (interface{}, error) {
	conn := rp.Get()
	defer conn.Close()
	v, err := conn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// GetMulti get redis from db.
func (rp *Pool) GetMulti(keys []string) ([]interface{}, error) {
	conn := rp.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("MGET", keys))
	if err != nil {
		return nil, err
	}
	return values, nil
}

// PutPrimitive Put put redis to db.
func (rp *Pool) PutPrimitive(key string, val interface{}) error {
	conn := rp.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, val)
	return err
}

// PutPrimitiveExpire Put put redis to db.
func (rp *Pool) PutPrimitiveExpire(key string, val interface{}, timeout time.Duration) error {
	conn := rp.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", key, int64(timeout/time.Second), val)
	return err
}

// GetStruct Get redis from db.
func (rp *Pool) GetStruct(key string, result interface{}) (exists bool, err error) {
	reply, err := rp.GetString(key)
	if err != nil || reply == "" {
		return false, err
	}

	err = msgpack.Unmarshal([]byte(reply), result)
	return true, err
}

// PutStruct Put put redis to db.
func (rp *Pool) PutStruct(key string, val interface{}) error {
	data, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}
	err = rp.PutPrimitive(key, string(data))
	return err
}

// PutStructExpire Put put redis to db.
func (rp *Pool) PutStructExpire(key string, val interface{}, timeout time.Duration) error {

	data, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}
	err = rp.PutPrimitiveExpire(key, string(data), timeout)
	return err
}

// Delete delete redis in db.
func (rp *Pool) Delete(key string) error {
	conn := rp.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

// IsExist check redis's existence in db.
func (rp *Pool) IsExist(key string) (bool, error) {
	conn := rp.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("EXISTS", key))
}

// Incr increase counter in db.
func (rp *Pool) Incr(key string, incrBy int) (int64, error) {
	conn := rp.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("INCRBY", key, incrBy))
}

// Decr decrease counter in db.
func (rp *Pool) Decr(key string) (int64, error) {
	conn := rp.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("DECR", key))
}

// Expire expire key
func (rp *Pool) Expire(key string, timeout time.Duration) error {
	conn := rp.Get()
	defer conn.Close()
	_, err := conn.Do("expire", key, int64(timeout/time.Second))
	return err
}
