package cache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	Decr(ctx context.Context, key string) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}

type RedisCache struct {
	opt   *cache.Options
	codec *cache.Cache
	ctx   context.Context
}

func (r *RedisCache) Has(key string) bool {
	return r.codec.Exists(r.ctx, key)
}

func (r *RedisCache) Get(key string, object interface{}) (Value, bool) {
	if r.shouldGetNative(object) {
		e := r.NativeCall(func(c RedisClient) error {
			cmd := c.Get(r.ctx, key)
			err := cmd.Err()
			if err == nil {
				return cmd.Scan(object)
			}
			return err
		})
		return Value{object}, e == nil
	} else {
		err := r.codec.Get(r.ctx, key, object)
		if err != nil {
			return Value{}, false
		}

		return Value{object}, true
	}
}

func (r *RedisCache) Put(key string, data interface{}, d time.Duration) error {
	if r.shouldNative(data) {
		return r.nativePut(key, data, d)
	} else {
		return r.codec.Set(&cache.Item{
			Key:   key,
			Value: data,
			TTL:   d,
			Ctx:   r.ctx,
		})
	}
}

func (r *RedisCache) nativePut(key string, data interface{}, d time.Duration) error {
	return r.NativeCall(func(c RedisClient) error {
		return c.Set(r.ctx, key, data, d).Err()
	})
}

func (r *RedisCache) shouldNative(data interface{}) bool {
	switch data.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

func (r *RedisCache) shouldGetNative(data interface{}) bool {
	switch data.(type) {
	case *int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64:
		return true
	}
	return false
}

func (r *RedisCache) Forever(key string, data interface{}) error {
	return r.Put(key, data, -1)
}

func (r *RedisCache) Remove(key string) bool {
	err := r.codec.Delete(r.ctx, key)

	return err == nil
}

func (r *RedisCache) Increment(key string) bool {
	if r.Has(key) {
		err := r.NativeCall(func(c RedisClient) error {
			return c.Incr(r.ctx, key).Err()
		})
		return err == nil
	}

	return false
}

func (r *RedisCache) NativeCall(f func(c RedisClient) error) error {
	if rc, ok := r.opt.Redis.(RedisClient); ok {
		return f(rc)
	}

	return errors.New("can not convert codec.Redis to RedisClient")
}

func (r *RedisCache) Decrement(key string) bool {
	if r.Has(key) {
		err := r.NativeCall(func(c RedisClient) error {
			return c.Decr(r.ctx, key).Err()
		})
		return err == nil
	}

	return false
}

func (r *RedisCache) Add(key string, data interface{}, d time.Duration) bool {
	if !r.Has(key) {
		if e := r.Put(key, data, d); e != nil {
			return false
		}
		return true
	}

	return false
}

func NewRedisCache(options *cache.Options) *RedisCache {
	ctx := context.Background()

	return &RedisCache{codec: cache.New(options), opt: options, ctx: ctx}
}
