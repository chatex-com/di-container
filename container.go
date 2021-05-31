package di

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const (
	tag          = "di"
	tagSkipField = "-"
	tagOptional  = "optional"
)

var (
	container               *Container
	ErrUnknownItem          = errors.New("unknown item")
	ErrElementMustBePointer = errors.New("element for loading must be pointer")
	ErrElementMustBeStruct  = errors.New("element for resolving must be pointer to struct")
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

	return c.load(&rv, name...)
}

func (c *Container) load(v *reflect.Value, name ...string) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	rv := *v
	rt := rv.Type()

	id := ""
	if len(name) > 0 {
		id = name[0]
	}

	id = c.generateID(rt, id)

	iv, exist := c.items[id]
	if !exist {
		return fmt.Errorf("%q: %w", id, ErrUnknownItem)
	}

	it := iv.Type()

	// pointer inside pointer case
	if rt.Kind() == reflect.Ptr && it.Kind() == reflect.Ptr {
		iv = iv.Elem()
		if !iv.CanAddr() {
			return fmt.Errorf("%q: %w", id, ErrUnableToSetPointer)
		}
		rv.Set(iv.Addr())

		return nil
	}

	if rt.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rt.Kind() != reflect.Interface && it.Kind() == reflect.Ptr {
		iv = iv.Elem()
	}

	if !rv.CanSet() {
		return fmt.Errorf("%q: %w", id, ErrUnableToSetValue)
	}

	rv.Set(iv)

	return nil
}

func (c *Container) Resolve(element interface{}) error {
	rv := reflect.ValueOf(element)
	rt := rv.Type()

	if rt.Kind() != reflect.Ptr {
		return ErrElementMustBePointer
	}

	rv = rv.Elem()
	rt = rv.Type()

	if rt.Kind() != reflect.Struct {
		return ErrElementMustBeStruct
	}

	numFields := rt.NumField()
	for i := 0; i < numFields; i++ {
		fv := rv.Field(i)
		ft := rv.Type().Field(i)

		if !fv.CanSet() {
			continue
		}

		name, opts := parseTag(ft.Tag.Get(tag))
		if name == tagSkipField {
			continue
		}

		err := c.load(&fv, name)

		if errors.Is(err, ErrUnknownItem) && opts.Contains(tagOptional) {
			continue
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Container) generateID(t reflect.Type, name string) string {
	if name != "" {
		return name
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return fmt.Sprintf("%s:%s.%s", name, t.PkgPath(), t.Name())
}

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return strings.TrimSpace(tag[:idx]), tagOptions(tag[idx+1:])
	}
	return strings.TrimSpace(tag), ""
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
