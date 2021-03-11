package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	DefaultExpiration = time.Minute * 20
	CleanupInterval   = time.Minute * 10
	KeyPrefix         = ""
)

type DriverRegister func() Repository

type CacheAble interface {
	MarshalToCache() interface{}
	UnmarshalFromCache(decoder func(value interface{}) bool) bool
}

var (
	DriverRegisters = make(map[string]DriverRegister)
	dm              = new(sync.RWMutex)
)

type Repository interface {
	Has(key string) bool
	Get(key string, object interface{}) (Value, bool)
	Put(key string, data interface{}, d time.Duration) error
	Forever(key string, data interface{}) error
	Remove(key string) bool
	Increment(key string) bool
	Decrement(key string) bool
	Add(key string, data interface{}, d time.Duration) bool
}

type Manager struct {
	driver     Repository
	driverName string
}

func (m *Manager) Has(key string) bool {
	return m.driver.Has(RealKey(key))
}

func (m *Manager) Get(key string, object interface{}) (Value, bool) {
	key = RealKey(key)
	if c, is := object.(CacheAble); is {
		var v Value
		result := c.UnmarshalFromCache(func(value interface{}) bool {
			vv, ok := m.driver.Get(key, value)
			v = vv
			return ok
		})
		return v, result
	} else {
		return m.driver.Get(key, object)
	}

}

func (m *Manager) Put(key string, data interface{}, d time.Duration) error {
	key = RealKey(key)

	if c, ok := data.(CacheAble); ok {
		data := c.MarshalToCache()
		return m.driver.Put(key, data, d)
	} else {
		return m.driver.Put(key, data, d)
	}
}

func (m *Manager) Forever(key string, data interface{}) error {
	return m.driver.Forever(RealKey(key), data)
}

func (m *Manager) Remove(key string) bool {
	return m.driver.Remove(RealKey(key))
}

func (m *Manager) Increment(key string) bool {
	return m.driver.Increment(RealKey(key))
}

func (m *Manager) Decrement(key string) bool {
	return m.driver.Decrement(RealKey(key))
}

func (m *Manager) Add(key string, data interface{}, d time.Duration) bool {
	return m.driver.Add(RealKey(key), data, d)
}

func (m *Manager) Use(driver string) error {
	if register, ok := getDriverRegister(driver); ok {
		if driver != m.driverName {
			m.driver = register()
			m.driverName = driver
		}
		return nil
	}

	return errors.New(fmt.Sprintf("cache: driver [%s] not registerd", driver))
}

func getDriverRegister(driver string) (DriverRegister, bool) {
	dm.RLock()
	defer dm.RUnlock()
	register, ok := DriverRegisters[driver]
	return register, ok
}

func RegisterDriver(name string, register DriverRegister) {
	dm.Lock()
	DriverRegisters[name] = register
	dm.Unlock()
}

type Value struct {
	d interface{}
}

func (v *Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.d)
}

func (v *Value) Data() interface{} {
	return v.d
}

func NewManager(defaultDriver ...string) *Manager {
	m := new(Manager)
	if len(defaultDriver) > 0 {
		m.Use(defaultDriver[0])
	}

	return m
}

func RealKey(k string) string {
	return KeyPrefix + k
}
