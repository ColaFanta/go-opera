package opera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsChannelNil(t *testing.T) {
	is := assert.New(t)

	var ch1 <-chan int

	ch2 := make(<-chan string)

	var ch3 chan<- bool

	is.True(isChannelNil[int](ch1))
	is.False(isChannelNil[string](ch2))
	is.True(isChannelNil[bool](ch3))
	is.True(isChannelNil[float64](nil))
}

func TestAwaitAnyReflect(t *testing.T) {
	is := assert.New(t)

	ctx := t.Context()

	ch1 := make(chan Result[int], 1)
	ch2 := make(chan Result[string], 1)

	ch1 <- Ok(42)
	ch2 <- Ok("Hello")

	res1 := awaitAnyReflect(ctx, ch1)
	is.True(res1.IsOk())
	is.Equal(42, res1.Yield())

	res2 := awaitAnyReflect(ctx, ch2)
	is.True(res2.IsOk())
	is.Equal("Hello", res2.Yield())
}
