package opera_test

import (
	"fmt"

	"github.com/colafanta/go-opera"
)

func ExampleTuple_Get() {
	name, age := opera.Tup("Alice", 30).Get()
	fmt.Printf("name: %s (%T)\n", name, name)
	fmt.Printf("age: %d (%T)\n", age, age)
	// Output:
	// name: Alice (string)
	// age: 30 (int)
}

func ExampleTuple3_Get() {
	name, age, active := opera.Tup3("Alice", 30, true).Get()
	fmt.Println(name, age, active)
	// Output: Alice 30 true
}
