package opera_test

import (
	"errors"
	"fmt"

	"github.com/colafanta/go-opera"
)

func ExampleOk() {
	status := "ready"
	result := opera.Ok(status)
	isReady := result.IsOk()
	fmt.Println("status:", result.OrEmpty())
	fmt.Println("isReady:", isReady)
	// Output:
	// status: ready
	// isReady: true
}

func ExampleErr() {
	validationErr := errors.New("invalid configuration")
	result := opera.Err[string](validationErr)
	returnedErr := result.Err()
	fmt.Println("error:", returnedErr)
	// Output: error: invalid configuration
}

func ExampleTry() {
	readPort := func() (int, error) { return 8080, nil }

	valOfPlainWay, plainErr := readPort()
	if plainErr != nil {
		fmt.Println("plain error:", plainErr)
		return
	}

	valOfOpera, operaErr := opera.Try(readPort()).Get()
	if operaErr != nil {
		fmt.Println("opera error:", operaErr)
		return
	}

	fmt.Println("valOfPlainWay:", valOfPlainWay)
	fmt.Println("valOfOpera:", valOfOpera)
	// Output:
	// valOfPlainWay: 8080
	// valOfOpera: 8080
}

func ExampleTryPass() {
	shutdownErr := errors.New("shutdown failed")
	result := opera.TryPass(shutdownErr)
	returnedErr := result.Err()
	fmt.Println("error:", returnedErr)
	// Output: error: shutdown failed
}

func ExampleTryHave() {
	status := "ready"
	result := opera.TryHave(status, true)
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleTryLookUp() {
	config := map[string]int{"port": 8080}

	valOfPlainWay, ok := config["port"]
	if !ok {
		fmt.Println("plain value is missing")
		return
	}

	result := opera.TryLookUp(config, "port")
	valOfOpera, operaErr := result.Get()
	if operaErr != nil {
		fmt.Println("opera error:", operaErr)
		return
	}

	fmt.Println("valOfPlainWay:", valOfPlainWay)
	fmt.Println("valOfOpera:", valOfOpera)
	// Output:
	// valOfPlainWay: 8080
	// valOfOpera: 8080
}

func ExampleTryAt() {
	stages := []string{"validate", "persist"}
	result := opera.TryAt(stages, 1)
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: persist
}

func ExampleTryIf() {
	isAuthorized := true
	result := opera.TryIf(isAuthorized)
	fmt.Println("isAuthorized:", result.IsOk())
	// Output: isAuthorized: true
}

func ExampleTryCast() {
	input := any("ready")
	result := opera.TryCast[string](input)
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleMustCast() {
	input := any("ready")
	result := opera.Do(func() string {
		castedStatus := opera.MustCast[string](input)
		return castedStatus
	})
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleMust() {
	readPort := func() (int, error) { return 8080, nil }
	result := opera.Do(func() int {
		port := opera.Must(readPort())
		return port
	})
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: 8080
}

func ExampleMustPass() {
	result := opera.Do(func() opera.Unit {
		opera.MustPass(nil)
		return opera.U
	})
	fmt.Println("completed:", result.IsOk())
	// Output: completed: true
}

func ExampleMustHave() {
	result := opera.Do(func() string {
		status := opera.MustHave("ready", true)
		return status
	})
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleMustLookUp() {
	config := map[string]int{"port": 8080}
	result := opera.Do(func() int {
		port := opera.MustLookUp(config, "port")
		return port
	})
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: 8080
}

func ExampleMustAt() {
	stages := []string{"validate", "persist"}
	result := opera.Do(func() string {
		stage := opera.MustAt(stages, 0)
		return stage
	})
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: validate
}

func ExampleMustTrue() {
	result := opera.Do(func() opera.Unit {
		opera.MustTrue(true)
		return opera.U
	})
	fmt.Println("completed:", result.IsOk())
	// Output: completed: true
}

func ExampleScoped() {
	result := opera.Scoped(
		func() opera.Result[string] { return opera.Ok("connection") },
		func(connection string) opera.Result[int] { return opera.Ok(len(connection)) },
		func(connection string) { fmt.Println("released", connection) },
	)
	fmt.Println(result.OrEmpty())
	// Output:
	// released connection
	// 10
}

func ExampleEmpty() {
	emptyName := opera.Empty[string]()
	fmt.Printf("emptyName: %q\n", emptyName)
	// Output: emptyName: ""
}

func ExampleIsEmpty() {
	emptyName := ""
	readyName := "ready"
	fmt.Println("emptyName:", opera.IsEmpty(emptyName))
	fmt.Println("readyName:", opera.IsEmpty(readyName))
	// Output:
	// emptyName: true
	// readyName: false
}

func ExampleDo() {
	readPort := func(config map[string]string) (int, error) {
		portText, ok := config["PORT"]
		if !ok {
			return 0, errors.New("PORT is missing")
		}
		if portText != "8080" {
			return 0, errors.New("invalid PORT")
		}
		return 8080, nil
	}
	config := map[string]string{"PORT": "8080"}

	valOfPlainWay, plainErr := readPort(config)
	if plainErr != nil {
		fmt.Println("plain error:", plainErr)
		return
	}

	result := opera.Do(func() int {
		port := opera.Try(readPort(config)).Yield()
		return port
	})
	valOfOpera, operaErr := result.Get()
	if operaErr != nil {
		fmt.Println("opera error:", operaErr)
		return
	}

	fmt.Println("valOfPlainWay:", valOfPlainWay)
	fmt.Println("valOfOpera:", valOfOpera)
	// Output:
	// valOfPlainWay: 8080
	// valOfOpera: 8080
}
