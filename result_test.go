package opera

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultOk(t *testing.T) {
	is := assert.New(t)

	is.Equal(Result[int]{val: 42, err: nil}, Ok(42))
}

func TestResultErr(t *testing.T) {
	is := assert.New(t)

	is.Nil(Ok(42).Err())
	is.NotNil(Err[int](assert.AnError).Err())
	is.Equal(assert.AnError, Err[int](assert.AnError).Err())
}

func TestResultTry(t *testing.T) {
	is := assert.New(t)

	is.Equal(Result[int]{val: 3, err: assert.AnError}, Try(3, assert.AnError))
}

func TestResultTryPass(t *testing.T) {
	is := assert.New(t)

	is.Equal(Result[Unit]{val: U, err: assert.AnError}, TryPass(assert.AnError))
}

func TestResultTryIf(t *testing.T) {
	is := assert.New(t)

	is.Equal(Ok(42), TryHave(42, true))
	is.Equal(Try(3, fmt.Errorf("%w: %T", ErrNoSuchElement, 3)), TryHave(3, false))
}

func TestResultTryAt(t *testing.T) {
	is := assert.New(t)

	slice := []int{10, 20, 30}
	array := [3]int{40, 50, 60}

	is.Equal(Ok(20), TryAt(slice, 1))
	is.Equal(Ok(60), TryAt(array[:], 2))
	is.ErrorIs(TryAt(slice, -1).Err(), ErrNoSuchElement)
	is.ErrorIs(TryAt(slice, 3).Err(), ErrNoSuchElement)
	is.ErrorIs(TryAt([]int{}, 0).Err(), ErrNoSuchElement)
}

func TestResultTryCast(t *testing.T) {
	is := assert.New(t)

	is.Equal(Ok(42), TryCast[int](42))
	is.Equal(
		Err[int](fmt.Errorf("%w: expect %T, actual %T", ErrInputCastFailed, *new(int), "")),
		TryCast[int]("not an int"),
	)

}

func TestResultIsOk(t *testing.T) {
	is := assert.New(t)

	is.True(Ok(42).IsOk())
	is.False(Err[int](assert.AnError).IsOk())
}

func TestResultIsError(t *testing.T) {
	is := assert.New(t)

	is.False(Ok(42).IsErr())
	is.True(Err[int](assert.AnError).IsErr())
}

func TestResultGet(t *testing.T) {
	is := assert.New(t)

	v1, err1 := Ok(42).Get()
	v2, err2 := Err[int](assert.AnError).Get()

	is.Equal(42, v1)
	is.Nil(err1)
	is.Error(assert.AnError, err1)

	is.Equal(0, v2)
	is.NotNil(err2)
	is.Error(assert.AnError, err2)
}

func TestResultMustOk(t *testing.T) {
	is := assert.New(t)

	is.NotPanics(func() {
		Ok(42).Yield()
	})
	is.Panics(func() {
		Err[int](assert.AnError).Yield()
	})

	is.Equal(42, Ok(42).Yield())
}

func TestResultMustAt(t *testing.T) {
	is := assert.New(t)

	slice := []int{10, 20, 30}
	array := [3]int{40, 50, 60}

	is.Equal(20, MustAt(slice, 1))
	is.Equal(60, MustAt(array[:], 2))
	is.Panics(func() {
		MustAt(slice, -1)
	})
	is.Panics(func() {
		MustAt(slice, 3)
	})
}

func TestResultOr(t *testing.T) {
	is := assert.New(t)

	is.Equal(42, Ok(42).Or(21))
	is.Equal(21, Err[int](assert.AnError).Or(21))
}

func TestResultOrEmpty(t *testing.T) {
	is := assert.New(t)

	is.Equal(42, Ok(42).OrEmpty())
	is.Equal(0, Err[int](assert.AnError).OrEmpty())
}

func TestResultToPointer(t *testing.T) {
	is := assert.New(t)

	val := 42
	ptr := Ok(42).ToPointer()
	nilPtr := Err[int](assert.AnError).ToPointer()

	is.NotNil(ptr)
	is.Equal(&val, ptr)
	is.Nil(nilPtr)
}

func TestResultToOption(t *testing.T) {
	is := assert.New(t)

	some := Ok(42).ToOption()
	none := Err[int](assert.AnError).ToOption()

	is.Equal(Some(42), some)
	is.Equal(None[int](), none)
}

func TestResultMap(t *testing.T) {
	is := assert.New(t)

	opt1 := Ok("hello").Map(strings.ToUpper)
	opt2 := Err[string](assert.AnError).Map(func(s string) string {
		is.Fail("should not be called")
		return "42"
	})

	is.Equal(Ok("HELLO"), opt1)
	is.Equal(Err[string](assert.AnError), opt2)
}

func TestResultMapErrIs(t *testing.T) {
	is := assert.New(t)

	err1 := errors.New("different error")
	opt1 := Ok(21).MapErrIs(assert.AnError, err1)

	opt2 := Err[int](assert.AnError).MapErrIs(assert.AnError, err1)

	err3 := errors.New("another error")
	opt3 := Try(42, assert.AnError).MapErrIs(err1, err3)

	is.Equal(Ok(21), opt1)
	is.Equal(Err[int](err1), opt2)
	is.Equal(assert.AnError, opt3.Err())
	is.Equal(42, opt3.val)
}

func TestResultMapErrorIsAny(t *testing.T) {
	is := assert.New(t)

	err1 := errors.New("different error")
	opt1 := Ok(21).MapErrIs(ErrAny, err1)

	opt2 := Err[int](assert.AnError).MapErrIs(ErrAny, err1)

	is.Equal(Ok(21), opt1)
	is.Equal(Err[int](err1), opt2)
}

func TestResultMapErrorString(t *testing.T) {
	is := assert.New(t)

	err1 := errors.New("different error")
	opt1 := Ok(21).MapErrorString("not found", err1)

	opt2 := Err[int](errors.New("file not found")).MapErrorString("not found", err1)

	opt3 := Try(42, errors.New("file not found")).MapErrorString("not found", err1)

	is.Equal(Ok(21), opt1)
	is.Equal(Err[int](err1), opt2)
	is.Equal(err1, opt3.Err())
	is.Equal(42, opt3.val)
}

func TestResultMapErr(t *testing.T) {
	is := assert.New(t)

	opt1 := Ok(21).MapErr(func(err error) error {
		is.Fail("should not be called")
		return err
	})
	err2 := errors.New("different error")
	opt2 := Err[int](assert.AnError).MapErr(func(err error) error {
		return err2
	})

	opt3 := Result[int]{val: 10, err: assert.AnError}.MapErr(func(err error) error {
		return nil
	})

	opt4 := Result[int]{val: 10, err: assert.AnError}.MapErr(func(err error) error {
		return err2
	})

	is.Equal(Ok(21), opt1)
	is.Equal(Err[int](err2), opt2)
	is.Equal(Ok(10), opt3)
	is.Equal(10, opt4.val)
}

func TestResultChain(t *testing.T) {
	is := assert.New(t)

	opt1 := Ok(21).Chain(func(i int) Result[int] {
		return Ok(42)
	})
	opt2 := Err[int](assert.AnError).Chain(func(i int) Result[int] {
		is.Fail("should not be called")
		return Ok(42)
	})

	is.Equal(Ok(42), opt1)
	is.Equal(Err[int](assert.AnError), opt2)
}

func TestResultDo(t *testing.T) {
	is := assert.New(t)

	opt1 := Do(func() int {
		return 42
	})
	opt2 := Do(func() int {
		TryPass(errors.New("some error")).Yield()
		return 21
	})

	is.Equal(Result[int]{val: 42, err: nil}, opt1)
	is.NotNil(opt2.err)
}

func TestResultScoped(t *testing.T) {
	is := assert.New(t)

	resourceAcquired := false
	resourceReleased := false

	result := Scoped(
		func() Result[string] {
			resourceAcquired = true
			return Ok("resource")
		},
		func(res string) Result[int] {
			is.Equal("resource", res)
			return Ok(42)
		},
		func(res string) {
			is.Equal("resource", res)
			resourceReleased = true
		},
	)

	is.True(resourceAcquired)
	is.True(resourceReleased)
	is.Equal(Ok(42), result)
}

func TestResultScopedWithError(t *testing.T) {
	is := assert.New(t)

	resourceAcquired := false
	resourceReleased := false

	result := Scoped(
		func() Result[string] {
			resourceAcquired = true
			return Ok("resource")
		},

		func(res string) Result[int] {
			is.Equal("resource", res)
			return Err[int](errors.New("use error"))
		},
		func(res string) {
			is.Equal("resource", res)
			resourceReleased = true
		},
	)

	is.True(resourceAcquired)
	is.True(resourceReleased)
	is.Equal(Err[int](errors.New("use error")), result)

	is.EqualError(result.err, "use error")
}

func TestResultTypeConvert(t *testing.T) {
	is := assert.New(t)

	res0 := Ok[any]("hello")

	converted0 := CastResult[string](res0)
	is.Equal(Ok("hello"), converted0)

	res1 := Ok(42)

	converted1 := CastResult[string](res1)
	is.Equal(
		Err[string](fmt.Errorf("%w: expect string, actual int", ErrResultCastFailed)),
		converted1,
	)

	res2 := Ok("st").ToAny()
	is.Equal(Ok[any]("st"), res2)
	converted2 := CastResult[string](res2)
	is.Equal(Ok("st"), converted2)

}

func TestResultCatchIs(t *testing.T) {
	is := assert.New(t)

	originalErr := errors.New("original error")
	res1 := Err[int](originalErr).CatchIs(originalErr, 42)
	is.Equal(Ok(42), res1)

	otherErr := errors.New("other error")
	res2 := Err[int](originalErr).CatchIs(otherErr, 21)
	is.Equal(Err[int](originalErr), res2)

	res3 := Err[int](originalErr).CatchIs(ErrAny, 84)
	is.Equal(Ok(84), res3)

	res4 := Ok(7).CatchIs(originalErr, 100)
	is.Equal(Ok(7), res4)

	res5 := Try("has value error", assert.AnError)
	is.Equal(assert.AnError, res5.CatchIs(otherErr, "recovered").Err())
}

func TestResultCatchString(t *testing.T) {
	is := assert.New(t)

	originalErr := errors.New("file not found")
	res1 := Err[int](originalErr).CatchString("not found", 42)
	is.Equal(Ok(42), res1)

	otherErr := errors.New("permission denied")
	res2 := Err[int](otherErr).CatchString("not found", 21)
	is.Equal(Err[int](otherErr), res2)

	res3 := Ok(7).CatchString("not found", 100)
	is.Equal(Ok(7), res3)
}

func TestResultCatch(t *testing.T) {
	is := assert.New(t)

	originalErr := errors.New("original error")
	res1 := Err[int](originalErr).Catch(func(err error) Result[int] {
		return Ok(42)
	})
	is.Equal(Ok(42), res1)

	otherErr := errors.New("other error")
	res2 := Err[int](originalErr).Catch(func(err error) Result[int] {
		return Err[int](otherErr)
	})
	is.Equal(Err[int](otherErr), res2)

	res3 := Ok(7).Catch(func(err error) Result[int] {
		is.Fail("should not be called")
		return Ok(100)
	})
	is.Equal(Ok(7), res3)

	res4 := Result[int]{val: 10, err: originalErr}.Catch(func(err error) Result[int] {
		return Err[int](otherErr)
	})
	is.Equal(otherErr, res4.err)
	is.Equal(10, res4.val)
}

func TestResultToAny(t *testing.T) {
	is := assert.New(t)

	res1 := Ok(42)
	anyRes1 := res1.ToAny()
	is.Equal(Ok[any](42), anyRes1)

	res2 := Err[int](assert.AnError)
	anyRes2 := res2.ToAny()
	is.Equal(Try[any](0, assert.AnError), anyRes2)

	res3 := Try("hello", assert.AnError)
	anyRes3 := res3.ToAny()
	is.Equal(Try[any]("hello", assert.AnError), anyRes3)
}

func TestResultTap(t *testing.T) {
	is := assert.New(t)

	called := false
	val := 0

	Ok(42).Tap(func(v int) {
		called = true
		val = v
	})
	is.True(called)
	is.Equal(42, val)

	called = false
	Err[int](assert.AnError).Tap(func(v int) {
		called = true
	})
	is.False(called)

	called = false
	r := Result[int]{val: 10, err: assert.AnError}.Tap(func(v int) {
		is.Fail("should not be called")
		called = true
		val = v
	})
	is.False(called)
	is.Equal(10, r.val)
}

func TestResultTapErr(t *testing.T) {
	is := assert.New(t)

	called := false
	var capturedErr error

	Ok(42).TapErr(func(err error) {
		called = true
	})
	is.False(called)

	called = false
	Err[int](assert.AnError).TapErr(func(err error) {
		called = true
		capturedErr = err
	})
	is.True(called)
	is.Equal(assert.AnError, capturedErr)

	called = false
	capturedErr = nil
	r := Result[int]{val: 10, err: assert.AnError}.TapErr(func(err error) {
		called = true
		capturedErr = err
	})
	is.True(called)
	is.Equal(assert.AnError, capturedErr)
	is.Equal(10, r.val)
}
