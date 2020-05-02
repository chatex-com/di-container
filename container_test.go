package di

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestContainer(t *testing.T) {
	Convey("test creation new container", t, func() {
		c := NewContainer()

		So(c, ShouldHaveSameTypeAs, &Container{})
		So(c.items, ShouldBeEmpty)
	})

	Convey("test default container", t, func() {
		c := GetDefaultContainer()

		So(c, ShouldHaveSameTypeAs, &Container{})
		So(c, ShouldEqual, container)
	})
}

func TestContainerSet(t *testing.T) {
	Convey("test set scalar values", t, func() {
		container := NewContainer()

		a := 123
		b := "foobar"
		c := 25.457
		d := true

		container.Set(a)
		container.Set(b)
		container.Set(c)
		container.Set(d)

		So(container.items, ShouldHaveLength, 4)
		So(container.items[container.generateID(reflect.TypeOf(a), "")].Interface(), ShouldEqual, a)
		So(container.items[container.generateID(reflect.TypeOf(b), "")].Interface(), ShouldEqual, b)
		So(container.items[container.generateID(reflect.TypeOf(c), "")].Interface(), ShouldEqual, c)
		So(container.items[container.generateID(reflect.TypeOf(d), "")].Interface(), ShouldEqual, d)
	})

	Convey("test set function value", t, func() {
		container := NewContainer()

		a := func() int { return 1 }
		container.Set(a)

		So(container.items, ShouldHaveLength, 1)
		So(container.items[container.generateID(reflect.TypeOf(a), "")].Interface(), ShouldEqual, a)
	})

	Convey("test set struct value", t, func() {
		container := NewContainer()

		type A struct {
			value string
		}

		a := A{value: "foo bar"}
		container.Set(a)

		val := container.items[container.generateID(reflect.TypeOf(a), "")].Interface()
		So(container.items, ShouldHaveLength, 1)
		So(val, ShouldHaveSameTypeAs, a)
		So(val.(A).value, ShouldEqual, a.value)
	})

	Convey("test set pointer to struct value", t, func() {
		container := NewContainer()

		a := struct{ value string }{value: "foo bar"}
		container.Set(&a)

		So(container.items, ShouldHaveLength, 1)
		So(container.items[container.generateID(reflect.TypeOf(a), "")].Interface(), ShouldEqual, &a)
	})

	Convey("test set with custom name", t, func() {
		container := NewContainer()

		a := 123
		b := "foobar"
		c := 25.457

		container.Set(a)
		container.Set(b, "b")
		container.Set(c, "c")

		So(container.items, ShouldHaveLength, 3)
		So(container.items[container.generateID(reflect.TypeOf(a), "")].Interface(), ShouldEqual, a)
		So(container.items[container.generateID(reflect.TypeOf(b), "b")].Interface(), ShouldEqual, b)
		So(container.items[container.generateID(reflect.TypeOf(c), "c")].Interface(), ShouldEqual, c)
	})
}

func TestContainerLoad(t *testing.T) {
	Convey("load scalar values", t, func() {
		container := NewContainer()

		expectedA := 123
		expectedB := "foobar"
		expectedC := 25.457
		expectedD := true

		container.Set(expectedA)
		container.Set(expectedB)
		container.Set(expectedC)
		container.Set(expectedD)

		var actualA int
		var actualB string
		var actualC float64
		var actualD bool

		errA := container.Load(&actualA)
		errB := container.Load(&actualB)
		errC := container.Load(&actualC)
		errD := container.Load(&actualD)

		So(container.items, ShouldHaveLength, 4)
		So(errA, ShouldBeNil)
		So(actualA, ShouldEqual, expectedA)
		So(errB, ShouldBeNil)
		So(actualB, ShouldEqual, expectedB)
		So(errC, ShouldBeNil)
		So(actualC, ShouldEqual, expectedC)
		So(errD, ShouldBeNil)
		So(actualD, ShouldEqual, expectedD)
	})

	Convey("load func value", t, func() {
		container := NewContainer()

		actual := func() int { return 123 }
		container.Set(actual)

		var expected func() int
		err := container.Load(&expected)

		So(err, ShouldBeNil)
		So(actual, ShouldEqual, expected)
		So(actual(), ShouldEqual, 123)
	})

	Convey("load struct value", t, func() {
		container := NewContainer()

		type A struct{ value string }
		actual := A{value: "foo bar"}
		container.Set(actual)

		var expected A
		err := container.Load(&expected)

		So(err, ShouldBeNil)
		So(actual, ShouldHaveSameTypeAs, expected)
		So(actual, ShouldNotEqual, expected)
		So(actual.value, ShouldEqual, expected.value)
	})

	Convey("load struct value as pointer", t, func() {
		container := NewContainer()

		type A struct{ value string }
		actual := &A{value: "foo bar"}
		container.Set(actual)

		var expected *A
		err := container.Load(&expected)

		So(err, ShouldBeNil)
		So(actual, ShouldHaveSameTypeAs, expected)
		So(actual, ShouldEqual, expected)
	})

	Convey("load value by name", t, func() {
		container := NewContainer()

		type A struct{ value string }
		container.Set(A{value: "foo bar"})
		container.Set(A{value: "foo baz"}, "baz")

		var actual A
		err := container.Load(&actual)

		So(err, ShouldBeNil)
		So(actual, ShouldHaveSameTypeAs, A{})
		So(actual.value, ShouldEqual, "foo bar")

		err = container.Load(&actual, "baz")
		So(err, ShouldBeNil)
		So(actual, ShouldHaveSameTypeAs, A{})
		So(actual.value, ShouldEqual, "foo baz")
	})

	Convey("try to load without pointer", t, func() {
		container := NewContainer()

		var a int
		err := container.Load(a)

		So(err, ShouldBeError)
		So(err, ShouldEqual, ErrElementMustBePointer)
	})

	Convey("try to load unknown item", t, func() {
		container := NewContainer()

		var a int
		err := container.Load(&a)

		So(err, ShouldBeError)
		So(err, ShouldEqual, ErrUnknownItem)
	})

	Convey("try to load pointer to variable", t, func() {
		container := NewContainer()

		expected := 123
		container.Set(&expected)

		var actual int
		err := container.Load(&actual)

		So(err, ShouldBeNil)
		So(actual, ShouldEqual, expected)
	})

	Convey("try to load variable to pointer to pointer to variable", t, func() {
		container := NewContainer()

		expected := 123
		container.Set(expected)

		var actual *int
		err := container.Load(&actual)

		So(err, ShouldBeError)
		So(err, ShouldEqual, ErrUnableToSetValue)
	})
}

func TestContainerResolve(t *testing.T) {
	Convey("resolve struct", t, func() {
		c := NewContainer()

		type A struct{}
		type B struct {
			Foo int
			Bar string
			Baz float64
			A   *A
		}

		c.Set(123)
		c.Set("foo bar")
		c.Set(345.67)
		a := &A{}
		c.Set(a)

		b := B{}
		err := c.Resolve(&b)

		So(err, ShouldBeNil)
		So(b.Foo, ShouldEqual, 123)
		So(b.Bar, ShouldEqual, "foo bar")
		So(b.Baz, ShouldEqual, 345.67)
		So(b.A, ShouldEqual, a)
	})

	Convey("support skip tag", t, func() {
		c := NewContainer()

		c.Set(123)
		c.Set("foo bar")

		type B struct {
			Foo int
			Bar string `di:"-"`
		}

		b := B{}
		err := c.Resolve(&b)

		So(err, ShouldBeNil)
		So(b.Foo, ShouldEqual, 123)
		So(b.Bar, ShouldBeEmpty)
	})

	Convey("support di name in tag", t, func() {
		c := NewContainer()

		c.Set(123)
		c.Set(234, "foo")

		type B struct {
			Foo int `di:"foo"`
		}

		b := B{}
		err := c.Resolve(&b)

		So(err, ShouldBeNil)
		So(b.Foo, ShouldEqual, 234)
	})

	Convey("support di optional tag", t, func() {
		c := NewContainer()

		type B struct {
			Foo string `di:",optional"`
		}

		b := B{}
		err := c.Resolve(&b)

		So(err, ShouldBeNil)
		So(b.Foo, ShouldBeEmpty)

		c.Set("bar")
		err = c.Resolve(&b)

		So(err, ShouldBeNil)
		So(b.Foo, ShouldEqual, "bar")
	})
}
