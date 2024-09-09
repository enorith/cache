package cache_test

import (
	"testing"
	"time"

	"github.com/enorith/cache"
	cache2 "github.com/go-redis/cache/v9"
	gc "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
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

func TestGoCache(t *testing.T) {
	gc := getGc()
	var str string
	gc.Get("test", &str)
	t.Log("test:", str)

	type c struct {
		a string
	}

	var str2 c
	gc.Put("test 2", c{a: "test aaa"}, time.Minute)
	gc.Get("test 2", &str2)

	t.Log("test 2:", str2)
}

func TestManager(t *testing.T) {
	cache.RegisterDriver("redis", func() (cache.Repository, error) {
		return getRc(), nil
	})
	m := cache.NewManager("redis")

	m.Put("cache_test:m", "test", time.Minute)
}

func TestGetAny(t *testing.T) {
	type foo struct {
		Name string
	}
	rc := getRc()
	e := rc.Put("cache:test_any", foo{
		Name: "test123",
	}, time.Minute)
	if e != nil {
		t.Fatalf("error put redis cache %v", e)
	}

	var data foo

	_, ok := rc.Get("cache:test_any", &data)
	if !ok {
		t.Fatalf("error get redis cache any")
	}

	t.Logf("get %v", data)
}

func getRc() *cache.RedisCache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "localhost:6379",
		},
		DB: 1,
	})

	return cache.NewRedisCache(&cache2.Options{
		Redis:        ring,
		LocalCache:   cache2.NewTinyLFU(1000, time.Minute),
		StatsEnabled: false,
	}, "enorith:")
}

func getGc() *cache.GoCache {
	return cache.NewGoCache(gc.New(5*time.Minute, 5*time.Minute), "enorith:")
}
