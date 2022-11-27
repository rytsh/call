package call_test

import (
	"fmt"

	"github.com/rytsh/call"
)

func Example_simple() {
	// create registry
	reg := call.NewReg().
		AddFunction("hababam", func() string { return "hababam" })

	// call function
	returns, err := reg.Call("hababam")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(returns[0].(string))
	// Output:
	// hababam
}
