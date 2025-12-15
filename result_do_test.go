package opera

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo_Success(t *testing.T) {
	is := assert.New(t)

	result := Do(func() string {
		return "Hello, World!"
	})

	is.False(result.IsErr())
	is.Equal("Hello, World!", result.Yield())
}

func TestDo_Error(t *testing.T) {
	is := assert.New(t)

	result := Do(func() string {
		TryPass(errors.New("something went wrong")).Yield()
		return "This will not be reached"
	})

	is.True(result.IsErr())
	is.EqualError(result.Err(), "something went wrong")
}

func TestDo_ComplexSuccess(t *testing.T) {
	is := assert.New(t)

	bookRoom := func(params map[string]string) Result[[]string] {
		return Do(func() []string {
			values := validateBooking(params).Yield()
			booking := createBooking(values["guest"]).Yield()
			room := assignRoom(booking, values["roomType"]).Yield()
			return []string{booking, room}
		})
	}

	params := map[string]string{
		"guest":    "Foo Bar",
		"roomType": "Suite",
	}

	result := bookRoom(params)
	is.False(result.IsErr())
	is.Equal(
		[]string{
			"Booking Created for: Foo Bar",
			"Room Assigned: Suite for Booking Created for: Foo Bar",
		},
		result.Yield(),
	)
}

func TestDo_ComplexError(t *testing.T) {
	is := assert.New(t)

	bookRoom := func(params map[string]string) Result[[]string] {
		return Do(func() []string {
			values := validateBooking(params).Yield()
			booking := createBooking(values["guest"]).Yield()
			room := assignRoom(booking, values["roomType"]).Yield()
			return []string{booking, room}
		})
	}

	params := map[string]string{
		"guest":    "",
		"roomType": "Suite",
	}

	result := bookRoom(params)
	is.True(result.IsErr())
	is.EqualError(result.Err(), "validation failed")
}

func TestDo_NonErrorPanic(t *testing.T) {
	is := assert.New(t)

	is.PanicsWithError("unexpected panic", func() {
		result := Do(func() string {
			panic(errors.New("unexpected panic"))
		})
		result.Yield()
	})
}

func BenchmarkDoAndGet(b *testing.B) {
	b.Run("ErrFromGet", func(b *testing.B) {
		bookRoom := func(params map[string]string) Result[[]string] {
			values, err := validateBooking(params).Get()
			if err != nil {
				return Err[[]string](err)
			}
			booking, err := createBooking(values["guest"]).Get()
			if err != nil {
				return Err[[]string](err)
			}
			room, err := assignRoom(booking, values["roomType"]).Get()
			if err != nil {
				return Err[[]string](err)
			}
			return Ok([]string{booking, room})
		}
		params := map[string]string{
			"guest":    "Foo Bar",
			"roomType": "Suite",
		}
		for b.Loop() {
			bookRoom(params)
		}
	})

	b.Run("ErrFromDo", func(b *testing.B) {
		bookRoom := func(params map[string]string) Result[[]string] {
			return Do(func() []string {
				values := validateBooking(params).Yield()
				booking := createBooking(values["guest"]).Yield()
				room := assignRoom(booking, values["roomType"]).Yield()
				return []string{booking, room}
			})
		}
		params := map[string]string{
			"guest":    "Foo Bar",
			"roomType": "Suite",
		}
		for b.Loop() {
			bookRoom(params)
		}
	})
}

var validateBooking = func(params map[string]string) Result[map[string]string] {
	if params["guest"] != "" && params["roomType"] != "" {
		return Ok(params)
	}
	return Err[map[string]string](errors.New("validation failed"))
}

var createBooking = func(guest string) Result[string] {
	if guest != "" {
		return Ok("Booking Created for: " + guest)
	}
	return Err[string](errors.New("booking creation failed"))
}

var assignRoom = func(booking string, roomType string) Result[string] {
	if roomType != "" {
		return Ok("Room Assigned: " + roomType + " for " + booking)
	}
	return Err[string](errors.New("room assignment failed"))
}
