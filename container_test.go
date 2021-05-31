package di

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestContainer(t *testing.T) {
	Convey("new container", t, func() {
		c := NewContainer()

		Convey("should have correct type", func() {
			So(c, ShouldHaveSameTypeAs, &Container{})
		})
		Convey("should be empty", func() {
			So(c.items, ShouldBeEmpty)
		})
	})

	Convey("default container", t, func() {
		c := GetDefaultContainer()

		Convey("should have correct type", func() {
			So(c, ShouldHaveSameTypeAs, &Container{})
		})
		Convey("should be reference to the exactly default container", func() {
			So(c, ShouldEqual, container)
		})
	})
}

func TestContainerSet(t *testing.T) {
	Convey("set scalar values", t, func() {
		container := NewContainer()

		a := 123
		b := "foobar"
		c := 25.457
		d := true

		container.Set(a)
		container.Set(b)
		container.Set(c)
		container.Set(d)

		Convey("should correctly change container's length", func() {
			So(container.items, ShouldHaveLength, 4)
		})
		Convey("values should be placed in the container according to their types", func() {
			Convey("integer", func() {
				So(container.items[container.generateID(reflect.TypeOf(a), "")].Interface(), ShouldEqual, a)
			})
			Convey("string", func() {
				So(container.items[container.generateID(reflect.TypeOf(b), "")].Interface(), ShouldEqual, b)
			})
			Convey("float", func() {
				So(container.items[container.generateID(reflect.TypeOf(c), "")].Interface(), ShouldEqual, c)
			})
			Convey("bool", func() {
				So(container.items[container.generateID(reflect.TypeOf(d), "")].Interface(), ShouldEqual, d)
			})
		})
	})

	Convey("set function value", t, func() {
		container := NewContainer()

		a := func() int { return 1 }
		container.Set(a)

		Convey("should correctly change container's length", func() {
			So(container.items, ShouldHaveLength, 1)
		})
		Convey("value should be placed in the container according type", func() {
			So(container.items[container.generateID(reflect.TypeOf(a), "")].Interface(), ShouldEqual, a)
		})
	})

	Convey("set struct value", t, func() {
		container := NewContainer()

		type A struct {
			value string
		}

		a := A{value: "foo bar"}
		container.Set(a)

		Convey("should correctly change container's length", func() {
			So(container.items, ShouldHaveLength, 1)
		})
		Convey("value should be placed in the container according type", func() {
			val := container.items[container.generateID(reflect.TypeOf(a), "")].Interface()
			So(val, ShouldHaveSameTypeAs, a)

			Convey("with correct value", func() {
				So(val.(A).value, ShouldEqual, a.value)
			})
		})
	})

	Convey("set pointer to struct value", t, func() {
		container := NewContainer()

		a := struct{ value string }{value: "foo bar"}
		container.Set(&a)

		Convey("should correctly change container's length", func() {
			So(container.items, ShouldHaveLength, 1)
		})
		Convey("container's element should be correct reference to the source variable", func() {
			So(container.items[container.generateID(reflect.TypeOf(a), "")].Interface(), ShouldEqual, &a)
		})
	})

	Convey("set scalars with custom name", t, func() {
		container := NewContainer()

		a := 123
		b := "foobar"
		c := 25.457

		container.Set(a)
		container.Set(b, "b")
		container.Set(c, "c")

		Convey("should correctly change container's length", func() {
			So(container.items, ShouldHaveLength, 3)
		})
		Convey("values should be placed in the container according to their type and passed name", func() {
			Convey("value without custom name", func() {
				So(container.items[container.generateID(reflect.TypeOf(a), "")].Interface(), ShouldEqual, a)
			})
			Convey("value with name 'b'", func() {
				So(container.items[container.generateID(reflect.TypeOf(b), "b")].Interface(), ShouldEqual, b)
			})
			Convey("value with name 'c'", func() {
				So(container.items[container.generateID(reflect.TypeOf(c), "c")].Interface(), ShouldEqual, c)
			})
		})
	})
}

type testInterface interface {
	do()
}
type testImpl struct {}
func (testImpl) do() {}

type testPointerImpl struct {}
func (*testPointerImpl) do() {}

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

		Convey("container should have correct length", func() {
			So(container.items, ShouldHaveLength, 4)
		})

		Convey("int value", func() {
			var actual int
			err := container.Load(&actual)

			Convey("shouldn't raise error", func() {
				So(err, ShouldBeNil)
			})

			Convey("should have correct resolved value", func() {
				So(actual, ShouldEqual, expectedA)
			})
		})

		Convey("string value", func() {
			var actual string
			err := container.Load(&actual)

			Convey("shouldn't raise error", func() {
				So(err, ShouldBeNil)
			})

			Convey("should have correct resolved value", func() {
				So(actual, ShouldEqual, expectedB)
			})
		})

		Convey("float value", func() {
			var actual float64
			err := container.Load(&actual)

			Convey("shouldn't raise error", func() {
				So(err, ShouldBeNil)
			})

			Convey("should have correct resolved value", func() {
				So(actual, ShouldEqual, expectedC)
			})
		})

		Convey("bool value", func() {
			var actual bool
			err := container.Load(&actual)

			Convey("shouldn't raise error", func() {
				So(err, ShouldBeNil)
			})

			Convey("should have correct resolved value", func() {
				So(actual, ShouldEqual, expectedD)
			})
		})
	})

	Convey("load func value", t, func() {
		container := NewContainer()

		actual := func() int { return 123 }
		container.Set(actual)

		var expected func() int
		err := container.Load(&expected)

		Convey("shouldn't raise an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("func should be resolved", func() {
			So(actual, ShouldEqual, expected)

			Convey("and callable with expected result", func() {
				So(actual(), ShouldEqual, 123)
			})
		})
	})

	Convey("load struct value", t, func() {
		container := NewContainer()

		type A struct{ value string }
		actual := A{value: "foo bar"}
		container.Set(actual)

		var expected A
		err := container.Load(&expected)

		Convey("shouldn't raise and error", func() {
			So(err, ShouldBeNil)
		})

		Convey("should resolve correct type", func() {
			So(actual, ShouldHaveSameTypeAs, expected)

			Convey("and should be a copy (not reference)", func() {
				So(actual, ShouldNotEqual, expected)

				Convey("with expected fields' values", func() {
					So(actual.value, ShouldEqual, expected.value)
				})
			})
		})
	})

	Convey("load struct value reference", t, func() {
		container := NewContainer()

		type A struct{ value string }
		actual := &A{value: "foo bar"}
		container.Set(actual)

		var expected *A
		err := container.Load(&expected)

		Convey("shouldn't raise an error", func() {
			So(err, ShouldBeNil)
		})

		Convey("should be correct type resolved", func() {
			So(actual, ShouldHaveSameTypeAs, expected)

			Convey("and expected reference to the source", func() {
				So(actual, ShouldEqual, expected)
			})
		})
	})

	Convey("load value by name", t, func() {
		container := NewContainer()

		type A struct{ value string }

		Convey("when we set to container 2 similar values with alias and not", func() {
			container.Set(A{value: "foo bar"})
			container.Set(A{value: "foo baz"}, "baz")

			Convey("value without alias should be resolved correct", func() {
				var actual A
				err := container.Load(&actual)

				Convey("shouldn't raise an error", func() {
					So(err, ShouldBeNil)

					Convey("and should have correct type", func() {
						So(actual, ShouldHaveSameTypeAs, A{})

						Convey("with expected field's value", func() {
							So(actual.value, ShouldEqual, "foo bar")
						})
					})
				})
			})

			Convey("value with alias should be resolved correct", func() {
				var actual A
				err := container.Load(&actual, "baz")

				Convey("shouldn't raise an error", func() {
					So(err, ShouldBeNil)

					Convey("and should have correct type", func() {
						So(actual, ShouldHaveSameTypeAs, A{})

						Convey("with expected field's value", func() {
							So(actual.value, ShouldEqual, "foo baz")
						})
					})
				})
			})
		})
	})

	Convey("loading not to pointer variable", t, func() {
		container := NewContainer()

		var a int
		err := container.Load(a)

		Convey("should raise an error", func() {
			So(err, ShouldBeError)

			Convey("with correct message", func() {
				So(errors.Is(err, ErrElementMustBePointer), ShouldBeTrue)
			})
		})
	})

	Convey("loading unknown item", t, func() {
		container := NewContainer()

		var a int
		err := container.Load(&a)

		Convey("should raise an error", func() {
			So(err, ShouldBeError)

			Convey("with correct message", func() {
				So(errors.Is(err, ErrUnknownItem), ShouldBeTrue)
			})
		})
	})

	Convey("loading into reference variable", t, func() {
		container := NewContainer()

		expected := 123
		container.Set(&expected)

		var actual int
		err := container.Load(&actual)

		Convey("shouldn't raise an error", func() {
			So(err, ShouldBeNil)

			Convey("and resolve expected value", func() {
				So(actual, ShouldEqual, expected)
			})
		})
	})

	Convey("loading non pointer variable into pointer to pointer variable", t, func() {
		container := NewContainer()

		expected := 123
		container.Set(expected)

		var actual *int
		err := container.Load(&actual)

		Convey("should raise an error", func() {
			So(err, ShouldBeError)

			Convey("with correct message", func() {
				So(errors.Is(err, ErrUnableToSetValue), ShouldBeTrue)
			})
		})
	})

	Convey("loading struct to the interface variable", t, func() {
		container := NewContainer()

		a := testImpl{}
		container.Set(&a, "implementation")

		var actual testInterface
		err := container.Load(&actual, "implementation")

		Convey("shouldn't raise an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("should be correct reference to the variable", func() {
			So(actual, ShouldNotBeNil)
		})
	})

	Convey("loading reference to the interface variable", t, func() {
		container := NewContainer()

		a := testPointerImpl{}
		container.Set(&a, "implementation")

		var actual testInterface
		err := container.Load(&actual, "implementation")

		Convey("shouldn't raise an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("should be correct reference to the variable", func() {
			So(actual, ShouldNotBeNil)
		})
	})
}

func TestContainerResolve(t *testing.T) {
	Convey("resolving struct without tags", t, func() {
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

		Convey("shouldn't raise an error", func() {
			So(err, ShouldBeNil)
		})

		Convey("should correct resolve struct fields", func() {
			Convey("integer", func() {
				So(b.Foo, ShouldEqual, 123)
			})
			Convey("string", func() {
				So(b.Bar, ShouldEqual, "foo bar")
			})
			Convey("float", func() {
				So(b.Baz, ShouldEqual, 345.67)
			})
			Convey("struct", func() {
				So(b.A, ShouldEqual, a)
			})
		})
	})

	Convey("resolving struct with skip-tag", t, func() {
		c := NewContainer()

		c.Set(123)
		c.Set("foo bar")

		type B struct {
			Foo int
			Bar string `di:"-"`
		}

		b := B{}
		err := c.Resolve(&b)

		Convey("shouldn't raise an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("should resolve untagged fields", func() {
			So(b.Foo, ShouldEqual, 123)

			Convey("and skip tagged fields", func() {
				So(b.Bar, ShouldBeEmpty)
			})
		})
	})

	Convey("resolving struct with di tag for named field", t, func() {
		c := NewContainer()

		c.Set(123)
		c.Set(234, "foo")

		type B struct {
			Foo int `di:"foo"`
		}

		b := B{}
		err := c.Resolve(&b)

		Convey("shouldn't raise an error", func() {
			So(err, ShouldBeNil)
		})

		Convey("should correct resolve value", func() {
			So(b.Foo, ShouldEqual, 234)
		})
	})

	Convey("resolving struct with optional tag", t, func() {
		c := NewContainer()

		type B struct {
			Foo string `di:",optional"`
		}

		Convey("when container doesn't contain required element", func() {
			var b B
			err := c.Resolve(&b)

			Convey("shouldn't raise an error", func() {
				So(err, ShouldBeNil)

				Convey("and should skip tagged field", func() {
					So(b.Foo, ShouldBeEmpty)
				})
			})
		})

		Convey("when container contains required element", func() {
			var b B
			c.Set("bar")
			err := c.Resolve(&b)

			Convey("shouldn't raise an error", func() {
				So(err, ShouldBeNil)

				Convey("and should correct resolve tagged field", func() {
					So(b.Foo, ShouldEqual, "bar")
				})
			})
		})
	})

	Convey("resolving struct with field with type interface to struct", t, func() {
		c := NewContainer()

		c.Set(&testImpl{}, "foo")

		type B struct {
			Foo testInterface `di:"foo"`
		}

		b := B{}
		err := c.Resolve(&b)

		Convey("shouldn't raise an error", func() {
			So(err, ShouldBeNil)
		})

		Convey("should correct resolve value", func() {
			So(b.Foo, ShouldNotBeNil)
		})
	})

	Convey("resolving struct with field with type interface to pointer to struct", t, func() {
		c := NewContainer()

		c.Set(&testPointerImpl{}, "foo")

		type B struct {
			Foo testInterface `di:"foo"`
		}

		b := B{}
		err := c.Resolve(&b)

		Convey("shouldn't raise an error", func() {
			So(err, ShouldBeNil)
		})

		Convey("should correct resolve value", func() {
			So(b.Foo, ShouldNotBeNil)
		})
	})
}
