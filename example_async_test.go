package opera_test

import (
	"context"
	"fmt"

	"github.com/colafanta/go-opera"
)

func ExampleAwait() {
	ctx := context.Background()
	fetchStatus := func(context.Context) (string, error) {
		return "ready", nil
	}
	statusTask := opera.Async(ctx, func(ctx context.Context) opera.Result[string] {
		return opera.Try(fetchStatus(ctx))
	})
	result := opera.Await(ctx, statusTask)
	valOfOpera, operaErr := result.Get()
	if operaErr != nil {
		fmt.Println("opera error:", operaErr)
		return
	}
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleAwaitAll() {
	ctx := context.Background()
	loadPort := func(context.Context) (int, error) { return 8080, nil }
	loadStatus := func(context.Context) (string, error) { return "ready", nil }
	portTask := opera.Async(ctx, func(ctx context.Context) opera.Result[int] {
		return opera.Try(loadPort(ctx))
	})
	statusTask := opera.Async(ctx, func(ctx context.Context) opera.Result[string] {
		return opera.Try(loadStatus(ctx))
	})
	result := opera.AwaitAll[any](ctx, portTask, statusTask)
	values, operaErr := result.Get()
	if operaErr != nil {
		fmt.Println("opera error:", operaErr)
		return
	}
	fmt.Println("port:", values[0])
	fmt.Println("status:", values[1])
	// Output:
	// port: 8080
	// status: ready
}

func ExampleAwaitAllT() {
	ctx := context.Background()
	loadPrimaryPort := func(context.Context) (int, error) { return 8080, nil }
	loadMetricsPort := func(context.Context) (int, error) { return 9090, nil }
	primaryPortTask := opera.Async(ctx, func(ctx context.Context) opera.Result[int] {
		return opera.Try(loadPrimaryPort(ctx))
	})
	metricsPortTask := opera.Async(ctx, func(ctx context.Context) opera.Result[int] {
		return opera.Try(loadMetricsPort(ctx))
	})
	result := opera.AwaitAllT(ctx, primaryPortTask, metricsPortTask)
	ports, operaErr := result.Get()
	if operaErr != nil {
		fmt.Println("opera error:", operaErr)
		return
	}
	fmt.Println("primaryPort:", ports[0])
	fmt.Println("metricsPort:", ports[1])
	// Output:
	// primaryPort: 8080
	// metricsPort: 9090
}

func ExampleAwaitFirst() {
	ctx := context.Background()
	failedTask := make(chan opera.Result[string], 1)
	failedTask <- opera.Err[string](fmt.Errorf("unavailable"))
	readyTask := make(chan opera.Result[string], 1)
	readyTask <- opera.Ok("ready")
	result := opera.AwaitFirst[string](ctx, failedTask, readyTask)
	firstReady, operaErr := result.Get()
	if operaErr != nil {
		fmt.Println("opera error:", operaErr)
		return
	}
	fmt.Println("firstReadyValue:", firstReady.V0)
	fmt.Println("firstReadyIndex:", firstReady.V1)
	// Output:
	// firstReadyValue: ready
	// firstReadyIndex: 1
}

func ExampleAsync() {
	ctx := context.Background()
	fetchProfile := func(context.Context) (string, error) { return "profile ready", nil }
	fetchSettings := func(context.Context) (string, error) { return "settings ready", nil }

	type plainResult struct {
		value string
		err   error
	}
	profileChannel := make(chan plainResult, 1)
	settingsChannel := make(chan plainResult, 1)
	go func() {
		value, err := fetchProfile(ctx)
		profileChannel <- plainResult{value: value, err: err}
	}()
	go func() {
		value, err := fetchSettings(ctx)
		settingsChannel <- plainResult{value: value, err: err}
	}()
	plainProfile := <-profileChannel
	plainSettings := <-settingsChannel
	if plainProfile.err != nil || plainSettings.err != nil {
		fmt.Println("plain error")
		return
	}
	valOfPlainWay := plainProfile.value + ", " + plainSettings.value

	profileTask := opera.Async(ctx, func(ctx context.Context) opera.Result[string] {
		return opera.Try(fetchProfile(ctx))
	})
	settingsTask := opera.Async(ctx, func(ctx context.Context) opera.Result[string] {
		return opera.Try(fetchSettings(ctx))
	})
	combinedResult := opera.AwaitAllT(ctx, profileTask, settingsTask)
	values, operaErr := combinedResult.Get()
	if operaErr != nil {
		fmt.Println("opera error:", operaErr)
		return
	}
	valOfOpera := values[0] + ", " + values[1]

	fmt.Println("valOfPlainWay:", valOfPlainWay)
	fmt.Println("valOfOpera:", valOfOpera)
	// Output:
	// valOfPlainWay: profile ready, settings ready
	// valOfOpera: profile ready, settings ready
}
