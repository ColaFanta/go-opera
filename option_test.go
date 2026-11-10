package opera

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionToResultErrors(t *testing.T) {
	missing := errors.New("missing configuration")
	invalid := errors.New("invalid configuration")
	tests := []struct {
		name string
		errs []error
		want []error
	}{
		{"no errors", nil, []error{ErrNoSuchElement}},
		{"nil error", []error{nil}, []error{ErrNoSuchElement}},
		{"all nil", []error{nil, nil}, []error{ErrNoSuchElement}},
		{"custom error", []error{missing}, []error{missing}},
		{"leading nil", []error{nil, missing}, []error{missing}},
		{"trailing nil", []error{missing, nil}, []error{missing}},
		{"joined errors", []error{missing, nil, invalid}, []error{missing, invalid}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := assert.New(t)
			value, err := None[int]().ToResult(test.errs...).Get()
			is.Zero(value)
			for _, target := range test.want {
				is.ErrorIs(err, target)
			}
			if len(test.errs) > 0 && errors.Join(test.errs...) != nil {
				is.NotErrorIs(err, ErrNoSuchElement)
			}
			is.Contains(err.Error(), ": int")
			is.Equal(Ok(42), Some(42).ToResult(test.errs...))
			is.Equal(Ok(0), Some(0).ToResult(test.errs...))
		})
	}
}

func TestOptionSome(t *testing.T) {
	is := assert.New(t)

	is.Equal(Option[int]{val: 42, hasVal: true}, Some(42))
}

func TestOptionNone(t *testing.T) {
	is := assert.New(t)

	is.Equal(Option[int]{hasVal: false}, None[int]())
}

func TestFromMayHave(t *testing.T) {
	is := assert.New(t)

	cb := func(v int, ok bool) func() (int, bool) {
		return func() (int, bool) {
			return v, ok
		}
	}

	is.Equal(Option[int]{hasVal: false}, MayHave(cb(42, false)()))
	is.Equal(Option[int]{hasVal: true, val: 42}, MayHave(cb(42, true)()))
}

func TestFromMayAt(t *testing.T) {
	is := assert.New(t)

	slice := []int{10, 20, 30}
	array := [3]int{40, 50, 60}

	is.Equal(Option[int]{hasVal: true, val: 20}, MayAt(slice, 1))
	is.Equal(Option[int]{hasVal: true, val: 60}, MayAt(array[:], 2))
	is.Equal(Option[int]{hasVal: false}, MayAt(slice, -1))
	is.Equal(Option[int]{hasVal: false}, MayAt(slice, 3))
	is.Equal(Option[int]{hasVal: false}, MayAt([]int{}, 0))
}

func TestOptionMaybeCast(t *testing.T) {
	is := assert.New(t)

	is.Equal(Option[int]{hasVal: true, val: 42}, MayCast[int](42))
	is.Equal(Option[int]{hasVal: false}, MayCast[int]("not an int"))
}

func TestOptionFromEmptyable(t *testing.T) {
	is := assert.New(t)

	is.Equal(Option[error]{hasVal: false}, MaybeEmpty[error](nil))
	is.Equal(Option[error]{hasVal: true, val: assert.AnError}, MaybeEmpty(assert.AnError))

	is.Equal(Option[int]{hasVal: false}, MaybeEmpty(0))
	is.Equal(Option[int]{hasVal: true, val: 42}, MaybeEmpty(42))

	is.Equal(Option[[3]string]{hasVal: false}, MaybeEmpty([3]string{}))
}

func TestOptionFromPointer(t *testing.T) {
	is := assert.New(t)

	is.Equal(Option[error]{hasVal: false}, MaybeNilPtr[error](nil))
	is.Equal(Option[error]{hasVal: true, val: assert.AnError}, MaybeNilPtr(&assert.AnError))

	zero := 0
	fortyTwo := 42
	is.Equal(Option[int]{hasVal: true, val: 0}, MaybeNilPtr(&zero))
	is.Equal(Option[int]{hasVal: true, val: 42}, MaybeNilPtr(&fortyTwo))
}

func TestOptionIsSome(t *testing.T) {
	is := assert.New(t)

	is.True(Some(42).IsSome())
	is.False(None[int]().IsSome())
}
func TestOptionIsNone(t *testing.T) {
	is := assert.New(t)

	is.False(Some(42).IsNone())
	is.True(None[int]().IsNone())
}

func TestOptionGet(t *testing.T) {
	is := assert.New(t)

	v1, ok1 := Some(42).Get()
	v2, ok2 := None[int]().Get()
	v3, ok3 := Option[string]{val: "xx", hasVal: false}.Get()

	is.Equal(42, v1)
	is.Equal(true, ok1)
	is.Equal(0, v2)
	is.Equal(false, ok2)
	is.Equal("", v3)
	is.Equal(false, ok3)
}

func TestOptionMustGet(t *testing.T) {
	is := assert.New(t)

	is.NotPanics(func() {
		Some(42).Yield()
	})
	is.Panics(func() {
		None[int]().Yield()
	})

	is.Equal(42, Some(42).Yield())
}

func TestOptionOr(t *testing.T) {
	is := assert.New(t)

	is.Equal(42, Some(42).Or(21))
	is.Equal(21, None[int]().Or(21))
}

func TestOptionOrEmpty(t *testing.T) {
	is := assert.New(t)

	is.Equal(42, Some(42).OrEmpty())
	is.Equal(0, None[int]().OrEmpty())
}

func TestOptionToPointer(t *testing.T) {
	is := assert.New(t)

	p := Some(42).ToPointer()
	is.NotNil(p)
	is.Equal(42, *p)

	is.Nil(None[int]().ToPointer())
}

func TestOptionMap(t *testing.T) {
	is := assert.New(t)

	opt1 := Some(21).Map(func(i int) int {
		return i * 2
	})
	opt2 := None[int]().Map(func(i int) int {
		is.Fail("should not be called")
		return 42
	})
	opt3 := Option[string]{val: "hello", hasVal: false}.Map(func(value string) string {
		is.Fail("should not be called")
		return "world"
	})

	is.Equal(Some(42), opt1)
	is.Equal(None[int](), opt2)
	is.Equal(None[string](), opt3)
}

func TestOptionMapNone(t *testing.T) {
	is := assert.New(t)

	opt1 := Some(21).OrElse(func() Option[int] {
		is.Fail("should not be called")
		return Some(42)
	})
	opt2 := None[int]().OrElse(func() Option[int] {
		return Some(42)
	})

	is.Equal(Some(21), opt1)
	is.Equal(Some(42), opt2)
}

func TestOptionFlatMap(t *testing.T) {
	is := assert.New(t)

	opt1 := Some(21).FlatMap(func(i int) Option[int] {
		return Some(42)
	})
	opt2 := None[int]().FlatMap(func(i int) Option[int] {
		return Some(42)
	})
	opt3 := Option[string]{val: "hello", hasVal: false}.FlatMap(func(s string) Option[string] {
		is.Fail("should not be called")
		return Some(strings.ToUpper(s))
	})

	is.Equal(Some(42), opt1)
	is.Equal(None[int](), opt2)
	is.Equal(None[string](), opt3)
}

func TestOptionTypeConvert(t *testing.T) {
	is := assert.New(t)

	opt1 := Some(42).Cast[int]()
	opt2 := Some("not an int").Cast[int]()
	opt3 := Some(42).Cast[string]()
	opt4 := None[int]().Cast[int]()

	optany := Some(42).Cast[any]()

	is.Equal(Some(42), opt1)
	is.Equal(None[int](), opt2)
	is.Equal(None[string](), opt3)
	is.Equal(None[int](), opt4)
	is.Equal(optany, Some[any](42))
	is.Equal(optany.Cast[int](), Some(42))
}

func TestOptionTap(t *testing.T) {
	is := assert.New(t)

	called := false
	val := 0
	opt1 := Some(21).Tap(func(i int) {
		called = true
		val = i * 2
	})
	is.True(called)
	is.Equal(42, val)
	is.Equal(Some(21), opt1)

	called = false
	opt2 := None[int]().Tap(func(i int) {
		called = true
	})
	is.False(called)
	is.Equal(None[int](), opt2)

	called = false
	opt3 := Option[string]{val: "hello", hasVal: false}.Tap(func(s string) {
		is.Fail("should not be called")
	})
	is.False(called)
	is.Equal("hello", opt3.val)
}
