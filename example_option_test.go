package opera_test

import (
	"errors"
	"fmt"

	"github.com/colafanta/go-opera"
)

func ExampleOption_IsSome() {
	configuredStatus := opera.Some("ready")
	missingStatus := opera.None[string]()
	fmt.Println("configuredStatus:", configuredStatus.IsSome())
	fmt.Println("missingStatus:", missingStatus.IsSome())
	// Output:
	// configuredStatus: true
	// missingStatus: false
}

func ExampleOption_IsNone() {
	configuredStatus := opera.Some("ready")
	missingStatus := opera.None[string]()
	fmt.Println("configuredStatus:", configuredStatus.IsNone())
	fmt.Println("missingStatus:", missingStatus.IsNone())
	// Output:
	// configuredStatus: false
	// missingStatus: true
}

func ExampleOption_Get() {
	statusOption := opera.Some("ready")
	valOfOpera, ok := statusOption.Get()
	if !ok {
		fmt.Println("status is missing")
		return
	}
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleOption_Yield() {
	missingStatus := opera.None[string]()
	result := opera.Do(func() string {
		return missingStatus.Yield()
	})
	returnedErr := result.Err()
	fmt.Println("error:", returnedErr)
	// Output: error: no such element
}

func ExampleOption_Or() {
	configuredStatus := opera.Some("ready")
	missingStatus := opera.None[string]()
	configuredValue := configuredStatus.Or("offline")
	fallbackValue := missingStatus.Or("offline")
	fmt.Println("configuredValue:", configuredValue)
	fmt.Println("fallbackValue:", fallbackValue)
	// Output:
	// configuredValue: ready
	// fallbackValue: offline
}

func ExampleOption_OrEmpty() {
	configuredStatus := opera.Some("ready")
	missingStatus := opera.None[string]()
	configuredValue := configuredStatus.OrEmpty()
	missingValue := missingStatus.OrEmpty()
	fmt.Println("configuredValue:", configuredValue)
	fmt.Printf("missingValue: %q\n", missingValue)
	// Output:
	// configuredValue: ready
	// missingValue: ""
}

func ExampleOption_OrPanic() {
	var recovered error
	func() {
		defer func() { recovered, _ = recover().(error) }()
		opera.None[string]().OrPanic()
	}()
	fmt.Println("recovered:", recovered)
	// Output: recovered: no such element
}

func ExampleOption_Map() {
	statusOption := opera.Some("ready")
	lengthOption := statusOption.Map(func(value string) int { return len(value) })
	valOfOpera := lengthOption.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: 5
}

func ExampleOption_FlatMap() {
	result := opera.Some("ready").FlatMap(func(value string) opera.Option[int] {
		return opera.Some(len(value))
	})
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: 5
}

func ExampleOption_Tap() {
	statusOption := opera.Some("ready")
	result := statusOption.Tap(func(value string) { fmt.Println("observed:", value) })
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output:
	// observed: ready
	// valOfOpera: ready
}

func ExampleOption_TapNone() {
	missingStatus := opera.None[string]()
	result := missingStatus.TapNone(func() { fmt.Println("using fallback") })
	isMissing := result.IsNone()
	fmt.Println("isMissing:", isMissing)
	// Output:
	// using fallback
	// isMissing: true
}

func ExampleOption_OrElse() {
	missingStatus := opera.None[string]()
	result := missingStatus.OrElse(func() opera.Option[string] { return opera.Some("offline") })
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: offline
}

func ExampleOption_ToResult() {
	statusOption := opera.Some("ready")
	result := statusOption.ToResult()
	valOfOpera, operaErr := result.Get()
	if operaErr != nil {
		fmt.Println("opera error:", operaErr)
		return
	}
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleOption_Cast() {
	result := opera.Some[any]("ready").Cast[string]()
	valOfOpera := result.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleOption_ToPointer() {
	statusOption := opera.Some("ready")
	statusPointer := statusOption.ToPointer()
	missingStatus := opera.None[string]()
	missingPointer := missingStatus.ToPointer()
	fmt.Println("pointedValue:", *statusPointer)
	fmt.Println("missingPointer:", missingPointer == nil)
	// Output:
	// pointedValue: ready
	// missingPointer: true
}

func ExampleOption_IsZero() {
	missingStatus := opera.None[string]()
	configuredStatus := opera.Some("ready")
	fmt.Println("missingStatus:", missingStatus.IsZero())
	fmt.Println("configuredStatus:", configuredStatus.IsZero())
	// Output:
	// missingStatus: true
	// configuredStatus: false
}

func ExampleSome() {
	status := "ready"
	statusOption := opera.Some(status)
	valOfOpera := statusOption.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleNone() {
	missingStatus := opera.None[string]()
	fmt.Println("isMissing:", missingStatus.IsNone())
	// Output: isMissing: true
}

func ExampleMayHave() {
	statusOption := opera.MayHave("ready", true)
	valOfOpera := statusOption.Or("offline")
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleMayLookUp() {
	config := map[string]int{"port": 8080}

	valOfPlainWay, ok := config["port"]
	if !ok {
		fmt.Println("plain value is missing")
		return
	}

	portOption := opera.MayLookUp(config, "port")
	valOfOpera, operaOK := portOption.Get()
	if !operaOK {
		fmt.Println("opera value is missing")
		return
	}

	fmt.Println("valOfPlainWay:", valOfPlainWay)
	fmt.Println("valOfOpera:", valOfOpera)
	// Output:
	// valOfPlainWay: 8080
	// valOfOpera: 8080
}

func ExampleMayAt() {
	stages := []string{"validate", "persist"}
	stageOption := opera.MayAt(stages, 0)
	valOfOpera := stageOption.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: validate
}

func ExampleMaybeEmpty() {
	emptyLabel := ""
	labelOption := opera.MaybeEmpty(emptyLabel)
	fmt.Println("isMissing:", labelOption.IsNone())
	// Output: isMissing: true
}

func ExampleMaybeNilPtr() {
	var labelPointer *string
	labelOption := opera.MaybeNilPtr(labelPointer)
	fmt.Println("isMissing:", labelOption.IsNone())
	// Output: isMissing: true
}

func ExampleMayCast() {
	input := any("ready")
	statusOption := opera.MayCast[string](input)
	valOfOpera := statusOption.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleOption_ToResult_error() {
	missingStatus := opera.None[string]()
	result := missingStatus.ToResult()
	_, returnedErr := result.Get()
	isMissing := errors.Is(returnedErr, opera.ErrNoSuchElement)
	fmt.Println("isMissing:", isMissing)
	// Output: isMissing: true
}

func ExampleOption_MarshalJSON() {
	statusOption := opera.Some("ready")
	encodedStatus, marshalErr := statusOption.MarshalJSON()
	if marshalErr != nil {
		fmt.Println("marshal error:", marshalErr)
		return
	}
	fmt.Println("encodedStatus:", string(encodedStatus))
	// Output: encodedStatus: "ready"
}

func ExampleOption_UnmarshalJSON() {
	var option opera.Option[string]
	encodedStatus := []byte(`"ready"`)
	unmarshalErr := option.UnmarshalJSON(encodedStatus)
	if unmarshalErr != nil {
		fmt.Println("unmarshal error:", unmarshalErr)
		return
	}
	valOfOpera := option.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleOption_MarshalText() {
	statusOption := opera.Some("ready")
	encodedStatus, marshalErr := statusOption.MarshalText()
	if marshalErr != nil {
		fmt.Println("marshal error:", marshalErr)
		return
	}
	fmt.Println("encodedStatus:", string(encodedStatus))
	// Output: encodedStatus: "ready"
}

func ExampleOption_UnmarshalText() {
	var option opera.Option[string]
	encodedStatus := []byte(`"ready"`)
	unmarshalErr := option.UnmarshalText(encodedStatus)
	if unmarshalErr != nil {
		fmt.Println("unmarshal error:", unmarshalErr)
		return
	}
	valOfOpera := option.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleOption_MarshalBinary() {
	statusOption := opera.Some("ready")
	data, marshalErr := statusOption.MarshalBinary()
	if marshalErr != nil {
		fmt.Println("marshal error:", marshalErr)
		return
	}
	var decoded opera.Option[string]
	decodeErr := decoded.UnmarshalBinary(data)
	if decodeErr != nil {
		fmt.Println("decode error:", decodeErr)
		return
	}
	valOfOpera := decoded.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleOption_UnmarshalBinary() {
	statusOption := opera.Some("ready")
	data, _ := statusOption.MarshalBinary()
	var option opera.Option[string]
	decodeErr := option.UnmarshalBinary(data)
	if decodeErr != nil {
		fmt.Println("decode error:", decodeErr)
		return
	}
	valOfOpera := option.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleOption_GobEncode() {
	statusOption := opera.Some("ready")
	data, encodeErr := statusOption.GobEncode()
	if encodeErr != nil {
		fmt.Println("encode error:", encodeErr)
		return
	}
	var decoded opera.Option[string]
	decodeErr := decoded.GobDecode(data)
	if decodeErr != nil {
		fmt.Println("decode error:", decodeErr)
		return
	}
	valOfOpera := decoded.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleOption_GobDecode() {
	statusOption := opera.Some("ready")
	data, _ := statusOption.GobEncode()
	var option opera.Option[string]
	decodeErr := option.GobDecode(data)
	if decodeErr != nil {
		fmt.Println("decode error:", decodeErr)
		return
	}
	valOfOpera := option.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: ready
}

func ExampleOption_Scan() {
	var option opera.Option[int]
	databaseValue := int64(7)
	scanErr := option.Scan(databaseValue)
	if scanErr != nil {
		fmt.Println("scan error:", scanErr)
		return
	}
	valOfOpera := option.OrEmpty()
	fmt.Println("valOfOpera:", valOfOpera)
	// Output: valOfOpera: 7
}

func ExampleOption_Value() {
	statusOption := opera.Some("ready")
	databaseValue, valueErr := statusOption.Value()
	if valueErr != nil {
		fmt.Println("value error:", valueErr)
		return
	}
	fmt.Println("databaseValue:", databaseValue)
	// Output: databaseValue: ready
}

func ExampleOption_Equal() {
	configuredStatus := opera.Some("ready")
	expectedStatus := opera.Some("ready")
	fmt.Println("isEqual:", configuredStatus.Equal(expectedStatus))
	// Output: isEqual: true
}
