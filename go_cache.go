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

	if object != nil {
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
	decoded := false

	switch t := to.(type) {
	case *int:
		*t = from.(int)
		decoded = true
	case *int8:
		*t = from.(int8)
		decoded = true
	case *int16:
		*t = from.(int16)
		decoded = true
	case *int64:
		*t = from.(int64)
		decoded = true
	case *uint:
		*t = from.(uint)
		decoded = true
	case *uint8:
		*t = from.(uint8)
		decoded = true
	case *uint16:
		*t = from.(uint16)
		decoded = true
	case *uint64:
		*t = from.(uint64)
		decoded = true
	case *string:
		*t = from.(string)
		decoded = true
	case *bool:
		*t = from.(bool)
		decoded = true
	case *float32:
		*t = from.(float32)
		decoded = true
	case *float64:
		*t = from.(float64)
		decoded = true
	}
	if !decoded {
		v := reflect.ValueOf(to)
		if v.Kind() == reflect.Ptr {
			dv := reflect.ValueOf(from)
			if dv.Kind() == reflect.Ptr {
				v.Elem().Set(dv.Elem())
				decoded = true
			} else {
				v.Elem().Set(dv)
				decoded = true
			}
		}
	}
	return decoded
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
