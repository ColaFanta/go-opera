package opera_test

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/colafanta/go-opera"
)

func ExampleResult_IsOk() {
	successfulPort := opera.Ok(8080)
	failedPort := opera.Err[int](errors.New("invalid port"))
	fmt.Println("successfulPort:", successfulPort.IsOk())
	fmt.Println("failedPort:", failedPort.IsOk())
	// Output:
	// successfulPort: true
	// failedPort: false
}

func ExampleResult_IsErr() {
	successfulPort := opera.Ok(8080)
	failedPort := opera.Err[int](errors.New("invalid port"))
	fmt.Println("successfulPort:", successfulPort.IsErr())
	fmt.Println("failedPort:", failedPort.IsErr())
	// Output:
	// successfulPort: false
	// failedPort: true
}

func ExampleResult_Yield() {
	failedPort := opera.Err[int](errors.New("invalid port"))
	result := opera.Do(func() int { return failedPort.Yield() })
	returnedErr := result.Err()
	fmt.Println("error:", returnedErr)
	// Output: error: invalid port
}

func ExampleResult_Err() {
	successfulPort := opera.Ok(8080)
	failedPort := opera.Err[int](errors.New("invalid port"))
	fmt.Println("successfulError:", successfulPort.Err())
	fmt.Println("failedError:", failedPort.Err())
	// Output:
	// successfulError: <nil>
	// failedError: invalid port
}

func ExampleResult_Get() {
	partialPort := opera.Try(8080, errors.New("health check failed"))
	value, returnedErr := partialPort.Get()
	fmt.Println("partialValue:", value)
	fmt.Println("error:", returnedErr)
	// Output:
	// partialValue: 8080
	// error: health check failed
}

func ExampleResult_Or() {
	successfulPort := opera.Ok(8080)
	failedPort := opera.Err[int](errors.New("invalid port"))
	fmt.Println("successfulValue:", successfulPort.Or(3000))
	fmt.Println("fallbackValue:", failedPort.Or(3000))
	// Output:
	// successfulValue: 8080
	// fallbackValue: 3000
}

func ExampleResult_OrF() {
	failedPort := opera.Err[int](errors.New("invalid"))
	fallbackPort := failedPort.OrF(func(err error) int { return len(err.Error()) * 1000 })
	fmt.Println("fallbackPort:", fallbackPort)
	// Output: fallbackPort: 7000
}

func ExampleResult_OrEmpty() {
	readyStatus := opera.Ok("ready")
	failedStatus := opera.Err[string](errors.New("invalid status"))
	fmt.Println("readyStatus:", readyStatus.OrEmpty())
	fmt.Printf("emptyStatus: %q\n", failedStatus.OrEmpty())
	// Output:
	// readyStatus: ready
	// emptyStatus: ""
}

func ExampleResult_OrPanic() {
	var recovered error
	func() {
		defer func() { recovered, _ = recover().(error) }()
		opera.Err[int](errors.New("invalid")).OrPanic()
	}()
	fmt.Println("recovered:", recovered)
	// Output: recovered: invalid
}

func ExampleResult_Map() {
	portResult := opera.Ok(8080)
	addressResult := portResult.Map(func(value int) string { return fmt.Sprintf("localhost:%d", value) })
	valOfOpera := addressResult.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: localhost:8080
}

func ExampleResult_FlatMap() {
	portText := opera.Ok("8080")
	result := portText.FlatMap(func(value string) opera.Result[int] {
		return opera.Ok(len(value))
	})
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: 4
}

func ExampleResult_Tap() {
	statusResult := opera.Ok("ready")
	result := statusResult.Tap(func(value string) { fmt.Println("observed:", value) })
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output:
	// observed: ready
	// valOfOpera: ready
}

func ExampleResult_MapErrIs() {
	notFound := errors.New("not found")
	failedPort := opera.Err[int](notFound)
	result := failedPort.MapErrIs(notFound, errors.New("service unavailable"))
	fmt.Println("mappedError:", result.Err())
	// Output: mappedError: service unavailable
}

func ExampleResult_MapErrorString() {
	failedRequest := opera.Err[int](errors.New("database unavailable"))
	result := failedRequest.MapErrorString("database", errors.New("retry request"))
	fmt.Println("mappedError:", result.Err())
	// Output: mappedError: retry request
}

func ExampleResult_MapErr() {
	failedRequest := opera.Err[int](errors.New("invalid request"))
	result := failedRequest.MapErr(func(err error) error {
		return fmt.Errorf("request: %w", err)
	})
	fmt.Println("mappedError:", result.Err())
	// Output: mappedError: request: invalid request
}

func ExampleResult_MapErrorAs() {
	readErr := &fs.PathError{Op: "open", Path: "config.json", Err: fs.ErrNotExist}
	partialConfig := opera.Try("partial", fmt.Errorf("load: %w", readErr))
	result := partialConfig.MapErrorAs(func(err *fs.PathError) error {
		return fmt.Errorf("configuration %s: %w", err.Path, err.Err)
	})
	value, err := result.Get()
	fmt.Println("partialValue:", value)
	fmt.Println("mappedError:", err)
	fmt.Println("isNotExist:", errors.Is(err, fs.ErrNotExist))
	// Output:
	// partialValue: partial
	// mappedError: configuration config.json: file does not exist
	// isNotExist: true
}

func ExampleResult_CatchAs() {
	readErr := &fs.PathError{Op: "open", Path: "config.json", Err: fs.ErrNotExist}
	failedConfig := opera.Err[string](fmt.Errorf("load: %w", readErr))
	value, err := failedConfig.CatchAs[*fs.PathError]("default configuration").Get()
	fmt.Println("value:", value)
	fmt.Println("error:", err)
	// Output:
	// value: default configuration
	// error: <nil>
}

func ExampleResult_Catch() {
	failedPort := opera.Err[int](errors.New("not found"))
	result := failedPort.Catch(func(err error) opera.Result[int] {
		return opera.Ok(80)
	})
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: 80
}

func ExampleResult_CatchIs() {
	notFound := errors.New("not found")
	failedPort := opera.Err[int](notFound)
	result := failedPort.CatchIs(notFound, 80)
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: 80
}

func ExampleResult_CatchString() {
	failedPort := opera.Err[int](errors.New("cache miss"))
	result := failedPort.CatchString("cache", 80)
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: 80
}

func ExampleResult_TapErr() {
	failedPort := opera.Err[int](errors.New("invalid port"))
	result := failedPort.TapErr(func(err error) { fmt.Println("logged:", err) })
	stillFailed := result.IsErr()
	fmt.Println("stillFailed:", stillFailed)
	// Output:
	// logged: invalid port
	// stillFailed: true
}

func ExampleResult_ToPointer() {
	statusResult := opera.Ok("ready")
	statusPointer := statusResult.ToPointer()
	failedStatus := opera.Err[string](errors.New("invalid status"))
	missingPointer := failedStatus.ToPointer()
	fmt.Println("pointedValue:", *statusPointer)
	fmt.Println("missingPointer:", missingPointer == nil)
	// Output:
	// pointedValue: ready
	// missingPointer: true
}

func ExampleResult_ToOption() {
	statusResult := opera.Ok("ready")
	statusOption := statusResult.ToOption()
	valOfOpera, ok := statusOption.Get()
	fmt.Println("valOfOpera:", valOfOpera)
	fmt.Println("isPresent:", ok)
	// Output:
	// valOfOpera: ready
	// isPresent: true
}

func ExampleResult_Cast() {
	inputResult := opera.Ok[any]("ready")
	result := inputResult.Cast[string]()
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}
