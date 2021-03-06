package cache_test

import (
	"github.com/enorith/cache"
	cache2 "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

func TestRedisCache_Put(t *testing.T) {
	rc := getRc()
	e := rc.Put("cache:test_int", 42, time.Minute)
	if e != nil {
		t.Fatalf("error put redis cache %v", e)
	}
	var v int
	rc.Get("cache:test_int", &v)
	if v != 42 {
		t.Fatalf("error get redis cache %d != 42", v)
	}
}

func TestRedisCache_Increment(t *testing.T) {
	rc := getRc()
	key := "cache:test_incr"
	e := rc.Put(key, 42, time.Minute)
	if e != nil {
		t.Fatalf("error put redis cache %v", e)
	}
	rc.Increment(key)
	var v int
	rc.Get(key, &v)
	if v != 43 {
		t.Fatalf("error Increment redis cache %d != 43", v)
	}
}

func getRc() *cache.RedisCache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "127.0.0.1:16379",
		},
	})

	return cache.NewRedisCache(&cache2.Options{
		Redis:        ring,
		LocalCache:   cache2.NewTinyLFU(1000, time.Minute),
		StatsEnabled: false,
	})
}
