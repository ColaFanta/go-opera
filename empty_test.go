package opera_test

import (
	"testing"

	"github.com/colafanta/go-opera"
	"github.com/stretchr/testify/assert"
)

func TestEmptyCompare(t *testing.T) {
	is := assert.New(t)
	var a1 string
	b1 := opera.Empty[string]()
	is.Equal(a1, b1)

	var a2 int
	b2 := opera.Empty[int]()
	is.Equal(a2, b2)

	var a3 [3]int
	b3 := opera.Empty[[3]int]()
	is.Equal(a3, b3)

}

func TestEmptyDetection(t *testing.T) {
	is := assert.New(t)

	is.True(opera.IsEmpty(""))
	is.True(opera.IsEmpty(0))
	is.True(opera.IsEmpty(0.0))
	is.True(opera.IsEmpty(false))
	is.True(opera.IsEmpty([]any{}))
	is.True(opera.IsEmpty(map[any]any{}))
	is.True(opera.IsEmpty(struct{}{}))
	is.True(opera.IsEmpty(zerod{Field: 0}))
	var emptyArr1 [3]int
	is.True(opera.IsEmpty(emptyArr1))
	var emptyUint1 uint32 = 0
	is.True(opera.IsEmpty(emptyUint1))

	is.False(opera.IsEmpty("hello"))
	is.False(opera.IsEmpty(42))
	is.False(opera.IsEmpty(3.14))
	is.False(opera.IsEmpty(true))
	is.False(opera.IsEmpty([]any{1, 2, 3}))
	is.False(opera.IsEmpty(map[any]any{"key": "value"}))
	is.False(opera.IsEmpty(struct{ Field int }{Field: 1}))
	is.False(opera.IsEmpty(zerod{Field: 1}))
	var nonEmptyArr1 [2]int = [2]int{1, 2}
	is.False(opera.IsEmpty(nonEmptyArr1))
	var nonEmptyUint1 uint64 = 10
	is.False(opera.IsEmpty(nonEmptyUint1))
}

type zerod struct {
	Field int
}

func (z zerod) IsZero() bool {
	return z.Field == 0
}
