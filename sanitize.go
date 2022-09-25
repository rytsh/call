package call

import (
	"reflect"
	"runtime"
	"strings"
)

// getFunctionName returns function name without package name.
//
// If function is anonymous, it will return "funcN" where N is number of function.
// Argument must be a function, otherwise it will panic.
func getFunctionName(fn reflect.Value) string {
	name := runtime.FuncForPC(fn.Pointer()).Name()
	// remove package name
	name = name[strings.LastIndex(name, ".")+1:]

	return name
}
