package call_test

import (
	"fmt"
	"reflect"
	"strings"

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

func Example_newOption() {
	myOption := call.OptionFunc{
		Name: "addExtraValue",
		Fn: func(args []reflect.Value, options ...string) ([]reflect.Value, error) {
			// do something with options
			return []reflect.Value{
				reflect.ValueOf(fmt.Sprintf("%s+%s", args[0].Interface(), strings.Join(options, "+"))),
			}, nil
		},
	}
	// create registry
	reg := call.NewReg(myOption).
		AddArgument("arg", "check").
		AddFunction("print", func(v string) string {
			return v
		})

	// call function
	returns, err := reg.CallWithArgs("print", "arg:addExtraValue=1,2,3")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(returns[0])
	// Output:
	// check+1+2+3
}
