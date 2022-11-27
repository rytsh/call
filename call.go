package call

import (
	"fmt"
	"reflect"
	"strings"
)

// Call calls function with name and uses already registered arguments.
func (r *Reg) Call(name string) ([]any, error) {
	return r.CallWithArgs(name, r.fn[name].Args...)
}

// CallWithArgs calls function with name and arguments.
func (r *Reg) CallWithArgs(name string, args ...string) ([]any, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	f, ok := r.fn[name]

	if !ok {
		return nil, fmt.Errorf("function %s not found", name)
	}

	fnArgs := make([]reflect.Value, 0)
	// get arguments
	for _, arg := range args {
		argPure := strings.SplitN(arg, r.GetDelimeter(), 2)[0]
		// parse argument options
		if v, ok := r.args[argPure]; ok {
			// do options
			vChanged, err := r.VisitOptions(arg, v)
			if err != nil {
				return nil, fmt.Errorf("failed VisitOption %w", err)
			}

			fnArgs = append(fnArgs, vChanged...)
		} else {
			return nil, fmt.Errorf("argument %s not found", arg)
		}
	}

	// check length is equal to function arguments
	if f.Fn.Type().IsVariadic() {
		if len(fnArgs) < f.Fn.Type().NumIn()-1 {
			return nil, fmt.Errorf("not enough arguments")
		}
	} else {
		if len(fnArgs) != f.Fn.Type().NumIn() {
			return nil, fmt.Errorf("argument count mismatch")
		}
	}

	// check last argument type with variadic
	if f.Fn.Type().IsVariadic() {
		fnArgType := f.Fn.Type().In(f.Fn.Type().NumIn() - 1).Elem()
		for i := f.Fn.Type().NumIn() - 1; i < len(fnArgs); i++ {
			argType := fnArgs[i].Type()
			if !argType.AssignableTo(fnArgType) {
				return nil, fmt.Errorf("variadic function: index %d argument %s type mismatch with function %s type", i, argType, fnArgType)
			}
		}
	} else if f.Fn.Type().NumIn() > 0 {
		fnArgType := f.Fn.Type().In(f.Fn.Type().NumIn() - 1)
		argType := fnArgs[len(fnArgs)-1].Type()
		if !argType.AssignableTo(fnArgType) {
			return nil, fmt.Errorf("function: index %d argument %s type mismatch with function %s type", len(fnArgs)-1, argType, fnArgType)
		}
	}

	// call function
	returnV := f.Fn.Call(fnArgs)

	// convert return values to []any
	returns := make([]any, len(returnV))
	for i, v := range returnV {
		returns[i] = v.Interface()
	}

	return returns, nil
}
