package di

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	container               *Container
	ErrUnknownItem          = errors.New("unknown item")
	ErrElementMustBePointer = errors.New("element for loading must be pointer")
	ErrUnableToSetPointer   = errors.New("unable to set pointer to element")
	ErrUnableToSetValue     = errors.New("unable to set value to the pointer")
)

type Container struct {
	items map[string]reflect.Value
	lock  sync.RWMutex
}

func GetDefaultContainer() *Container {
	if container == nil {
		container = NewContainer()
	}

	return container
}

func NewContainer() *Container {
	return &Container{
		items: make(map[string]reflect.Value),
	}
}

func (c *Container) Set(element interface{}, name ...string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	id := ""
	if len(name) > 0 {
		id = name[0]
	}

	v := reflect.ValueOf(element)

	id = c.generateID(v.Type(), id)

	c.items[id] = v
}

func (c *Container) Load(element interface{}, name ...string) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	rv := reflect.ValueOf(element)

	if rv.Type().Kind() != reflect.Ptr {
		return ErrElementMustBePointer
	}

	rv = rv.Elem()
	rt := rv.Type()

	id := ""
	if len(name) > 0 {
		id = name[0]
	}

	id = c.generateID(rt, id)

	iv, exist := c.items[id]
	if !exist {
		return ErrUnknownItem
	}

	it := iv.Type()

	// pointer inside pointer case
	if rt.Kind() == reflect.Ptr && it.Kind() == reflect.Ptr {
		iv = iv.Elem()
		if !iv.CanAddr() {
			return ErrUnableToSetPointer
		}
		rv.Set(iv.Addr())

		return nil
	}

	if rt.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if it.Kind() == reflect.Ptr {
		iv = iv.Elem()
	}

	if !rv.CanSet() {
		return ErrUnableToSetValue
	}

	rv.Set(iv)

	return nil
}

func (c *Container) generateID(t reflect.Type, name string) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return fmt.Sprintf("%s:%s.%s", name, t.PkgPath(), t.Name())
}
