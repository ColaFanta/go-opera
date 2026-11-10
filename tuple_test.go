package opera_test

import (
	"reflect"
	"testing"

	"github.com/colafanta/go-opera"
	"github.com/stretchr/testify/assert"
)

func TestTupleGet(t *testing.T) {
	want := []any{"first", 2, true, int64(4), uint(5), float32(6), float64(7), byte(8), rune(9)}
	tests := []struct {
		name  string
		tuple any
	}{
		{"Tuple", opera.Tup("first", 2)},
		{"Tuple3", opera.Tup3("first", 2, true)},
		{"Tuple4", opera.Tup4("first", 2, true, int64(4))},
		{"Tuple5", opera.Tup5("first", 2, true, int64(4), uint(5))},
		{"Tuple6", opera.Tup6("first", 2, true, int64(4), uint(5), float32(6))},
		{"Tuple7", opera.Tup7("first", 2, true, int64(4), uint(5), float32(6), float64(7))},
		{"Tuple8", opera.Tup8("first", 2, true, int64(4), uint(5), float32(6), float64(7), byte(8))},
		{"Tuple9", opera.Tup9("first", 2, true, int64(4), uint(5), float32(6), float64(7), byte(8), rune(9))},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := assert.New(t)
			tuple := reflect.ValueOf(test.tuple)
			values := tuple.MethodByName("Get").Call(nil)
			if !is.Len(values, tuple.NumField()) {
				return
			}
			for index, value := range values {
				is.Equal(want[index], value.Interface())
				is.Equal(reflect.TypeOf(want[index]), value.Type())
			}

			zero := reflect.Zero(tuple.Type())
			zeroValues := zero.MethodByName("Get").Call(nil)
			if !is.Len(zeroValues, tuple.NumField()) {
				return
			}
			for index, value := range zeroValues {
				is.Equal(reflect.TypeOf(want[index]), value.Type())
				is.True(value.IsZero())
			}
		})
	}
}
