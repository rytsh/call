package call_test

import (
	"fmt"

	"github.com/rytsh/call"
)

func Example_variadic() {
	// create registry
	reg := call.NewReg().
		AddArgument("a", 6).
		AddArgument("b", 2).
		AddFunction("sum", func(x ...int) int {
			sum := 0
			for _, v := range x {
				sum += v
			}

			return sum
		})

	// call function
	returns, err := reg.CallWithArgs("sum", "a", "b")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(returns[0])
	// Output:
	// 8
}

func Example_options() {
	// create registry
	reg := call.NewReg().
		AddArgument("a", []int{2, 5}).
		AddFunction("sum", func(x ...int) int {
			sum := 0
			for _, v := range x {
				sum += v
			}

			return sum
		})

	// call function
	returns, err := reg.CallWithArgs("sum", "a:...")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(returns[0])

	// call function with a's elements
	returns, err = reg.CallWithArgs("sum", "a:index=0,1")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(returns[0])

	reg.AddArgument("b", map[string]int{"a": 2, "b": 5})

	// call function with b's elements
	returns, err = reg.CallWithArgs("sum", "b:...")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(returns[0])

	// call function with b's element
	returns, err = reg.CallWithArgs("sum", "b:index=b")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(returns[0])
	// Output:
	// 7
	// 7
	// 7
	// 5
}
