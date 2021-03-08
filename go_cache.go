package cache

import (
	"reflect"
	"time"

	gc "github.com/patrickmn/go-cache"
)

type GoCache struct {
	gc *gc.Cache
}

func (c *GoCache) Has(key string) bool {
	_, exists := c.gc.Get(key)
	return exists
}

func (c *GoCache) Get(key string, object interface{}) (Value, bool) {
	data, exists := c.gc.Get(key)

	unmarshal(data, object)

	return Value{data}, exists
}

func (c *GoCache) Put(key string, data interface{}, d time.Duration) error {
	c.gc.Set(key, data, d)
	return nil
}

func (c *GoCache) Forever(key string, data interface{}) error {
	c.gc.Set(key, data, gc.NoExpiration)
	return nil
}

func (c *GoCache) Remove(key string) bool {
	c.gc.Delete(key)
	return !c.Has(key)
}

func (c *GoCache) Increment(key string) bool {
	err := c.gc.Increment(key, 1)

	return err != nil
}

func (c *GoCache) Decrement(key string) bool {
	err := c.gc.Decrement(key, 1)

	return err != nil
}

func (c *GoCache) Add(key string, data interface{}, d time.Duration) bool {
	err := c.gc.Add(key, data, d)

	return err != nil
}

func unmarshal(from, to interface{}) bool {

	decoded := false

	switch t := to.(type) {
	case *int:
		*t = from.(int)
		decoded = true
		break
	case *int8:
		*t = from.(int8)
		decoded = true
		break
	case *int16:
		*t = from.(int16)
		decoded = true
		break
	case *int64:
		*t = from.(int64)
		decoded = true
		break
	case *uint:
		*t = from.(uint)
		decoded = true
		break
	case *uint8:
		*t = from.(uint8)
		decoded = true
		break
	case *uint16:
		*t = from.(uint16)
		decoded = true
		break
	case *uint64:
		*t = from.(uint64)
		decoded = true
		break
	case *string:
		*t = from.(string)
		decoded = true
		break
	case *bool:
		*t = from.(bool)
		decoded = true
		break
	case *float32:
		*t = from.(float32)
		decoded = true
		break
	case *float64:
		*t = from.(float64)
		decoded = true
		break
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

func NewGoCache(gc *gc.Cache) *GoCache {
	return &GoCache{gc: gc}
}
