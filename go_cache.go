package cache

import (
	"reflect"
	"time"

	gc "github.com/patrickmn/go-cache"
)

type GoCache struct {
	gc     *gc.Cache
	prefix string
}

func (c *GoCache) Has(key string) bool {
	_, exists := c.gc.Get(c.resolveKey(key))
	return exists
}

func (c *GoCache) Get(key string, object interface{}) (Value, bool) {
	data, exists := c.gc.Get(c.resolveKey(key))

	if object != nil && data != nil {
		unmarshal(data, object)
	}

	return Value{data}, exists
}

func (c *GoCache) Put(key string, data interface{}, d time.Duration) error {
	c.gc.Set(c.resolveKey(key), data, d)
	return nil
}

func (c *GoCache) Forever(key string, data interface{}) error {
	c.gc.Set(c.resolveKey(key), data, gc.NoExpiration)
	return nil
}

func (c *GoCache) Remove(key string) bool {
	c.gc.Delete(c.resolveKey(key))
	return !c.Has(c.resolveKey(key))
}

func (c *GoCache) Increment(key string) bool {
	err := c.gc.Increment(c.resolveKey(key), 1)

	return err != nil
}

func (c *GoCache) Decrement(key string) bool {
	err := c.gc.Decrement(c.resolveKey(key), 1)

	return err != nil
}

func (c *GoCache) Add(key string, data interface{}, d time.Duration) bool {
	err := c.gc.Add(c.resolveKey(key), data, d)

	return err != nil
}

func unmarshal(from, to interface{}) bool {
	val := reflect.Indirect(reflect.ValueOf(from))

	if val.IsZero() {
		return false
	}

	toVal := reflect.ValueOf(to)

	reflect.Indirect(toVal).Set(val)

	return true
}

func (c *GoCache) resolveKey(key string) string {
	return c.prefix + key
}

func NewGoCache(gc *gc.Cache, prefix ...string) *GoCache {
	var p string
	if len(prefix) > 0 {
		p = prefix[0]
	}

	return &GoCache{gc: gc, prefix: p}
}
