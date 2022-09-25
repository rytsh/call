package call_test

import (
	"fmt"

	"github.com/rytsh/call"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("divide by zero")
	}

	return a / b, nil
}

func Example() {
	// create registry
	reg := call.NewReg().
		AddArgument("a", 6).
		AddArgument("b", 2).
		AddFunction("", divide, "a", "b")

	// call function
	returns, err := reg.Call("divide")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(returns[0], returns[1])
	// Output:
	// 3 <nil>
}
