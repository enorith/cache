package cache

import (
	"errors"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
)

type RedisClient interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Incr(key string) *redis.IntCmd
	Decr(key string) *redis.IntCmd
	Get(key string) *redis.StringCmd
}

type RedisCache struct {
	codec *cache.Codec
}

func (r *RedisCache) Has(key string) bool {
	return r.codec.Exists(key)
}

func (r *RedisCache) Get(key string, object interface{}) (Value, bool) {
	if r.shouldGetNative(object) {
		e := r.NativeCall(func(c RedisClient) error {
			cmd := c.Get(key)
			err := cmd.Err()
			if err == nil {
				return cmd.Scan(object)
			}
			return err
		})
		return Value{object}, e == nil
	} else {
		err := r.codec.Get(key, object)
		if err != nil {
			return Value{}, false
		}

		return Value{object}, true
	}
}

func (r *RedisCache) Put(key string, data interface{}, d time.Duration) {
	if r.shouldNative(data) {
		r.nativePut(key, data, d)
	} else {
		r.codec.Set(&cache.Item{
			Key:        key,
			Object:     data,
			Expiration: d,
		})
	}
}

func (r *RedisCache) nativePut(key string, data interface{}, d time.Duration) {
	r.NativeCall(func(c RedisClient) error {
		return c.Set(key, data, d).Err()
	})
}

func (r *RedisCache) shouldNative(data interface{}) bool {
	switch data.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	}
	return false
}

func (r *RedisCache) shouldGetNative(data interface{}) bool {
	switch data.(type) {
	case *int:
		return true
	case *int8:
		return true
	case *int16:
		return true
	case *int32:
		return true
	case *int64:
		return true
	case *uint:
		return true
	case *uint8:
		return true
	case *uint16:
		return true
	case *uint32:
		return true
	case *uint64:
		return true
	}
	return false
}

func (r *RedisCache) Forever(key string, data interface{}) {
	r.Put(key, data, -1)
}

func (r *RedisCache) Remove(key string) bool {
	err := r.codec.Delete(key)

	return err == nil
}

func (r *RedisCache) Increment(key string) bool {
	if r.Has(key) {
		err := r.NativeCall(func(c RedisClient) error {
			return c.Incr(key).Err()
		})
		return err == nil
	}

	return false
}

func (r *RedisCache) NativeCall(f func(c RedisClient) error) error {
	if rc, ok := r.codec.Redis.(RedisClient); ok {
		return f(rc)
	}

	return errors.New("can not convert codec.Redis to RedisClient")
}

func (r *RedisCache) Decrement(key string) bool {
	if r.Has(key) {
		err := r.NativeCall(func(c RedisClient) error {
			return c.Decr(key).Err()
		})
		return err == nil
	}

	return false
}

func (r *RedisCache) Add(key string, data interface{}, d time.Duration) bool {
	if !r.Has(key) {
		r.Put(key, data, d)
	}

	return false
}

func NewRedisCache(codec *cache.Codec) *RedisCache {
	return &RedisCache{codec: codec}
}
