package call

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Option interface {
	GetDelimeter() string
	AddOption(name string, fn func([]reflect.Value, ...string) ([]reflect.Value, error)) Option
	GetOption(name string) (func([]reflect.Value, ...string) ([]reflect.Value, error), bool)
	VisitOptions(arg string, v any) (error, []reflect.Value)
}

type Func struct {
	Args []string
	Fn   reflect.Value
}

// Reg is a registry for functions and arguments.
type Reg struct {
	fn    map[string]Func
	args  map[string]any
	mutex sync.RWMutex
	Option
}

// NewReg creates new registry.
func NewReg() *Reg {
	return &Reg{
		fn:   make(map[string]Func),
		args: make(map[string]any),
		Option: NewOptions().
			AddOption("index", OptionGetIndex).
			AddOption("...", OptionVariadic),
	}
}

// AddArgument adds argument to registry with name.
//
// If name includes delimeter, it will not add options.
func (r *Reg) AddArgument(name string, x any) *Reg {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// trim options
	name = strings.SplitN(name, r.GetDelimeter(), 2)[0]

	r.args[name] = x

	return r
}

// GetArgument returns argument with name.
func (r *Reg) GetArgument(name string) (any, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	v, ok := r.args[name]

	return v, ok
}

// DeleteArgument deletes argument with name.
func (r *Reg) DeleteArgument(name string) *Reg {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.args, name)

	return r
}

// GetArgumentNames returns all argument names.
func (r *Reg) GetArgumentNames() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.args))

	for name := range r.args {
		names = append(names, name)
	}

	return names
}

// AddFunction adds function to registry with name.
//
// If name is empty, function name will be used.
// Argument must be a function, otherwise it will panic.
func (r *Reg) AddFunction(name string, fn any, args ...string) *Reg {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	fnV := reflect.ValueOf(fn)
	if fnV.Kind() != reflect.Func {
		panic("fn is not a function")
	}

	if name == "" {
		name = getFunctionName(fnV)
	}

	r.fn[name] = Func{
		Args: args,
		Fn:   fnV,
	}

	return r
}

// GetFunction returns function with name.
func (r *Reg) GetFunction(name string) (Func, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	v, ok := r.fn[name]

	return v, ok
}

// RemoveFunction removes function with name.
func (r *Reg) RemoveFunction(name string) *Reg {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.fn, name)

	return r
}

// GetFunctionNames returns all function names.
func (r *Reg) GetFunctionNames() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.fn))

	for name := range r.fn {
		names = append(names, name)
	}

	return names
}

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
			err, vChanged := r.VisitOptions(arg, v)
			if err != nil {
				return nil, fmt.Errorf("error on VisitOption %w", err)
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

	// check arguments type with function arguments with variadic
	for i := 0; i < f.Fn.Type().NumIn()-2; i++ {
		fnArgType := f.Fn.Type().In(i)
		argType := fnArgs[i].Type()
		if !argType.AssignableTo(fnArgType) {
			return nil, fmt.Errorf("%d argument %s type mismatch with function %s type", i, argType, fnArgType)
		}
	}

	// check last argument type with variadic
	if f.Fn.Type().IsVariadic() {
		fnArgType := f.Fn.Type().In(f.Fn.Type().NumIn() - 1).Elem()
		for i := f.Fn.Type().NumIn() - 1; i < len(fnArgs); i++ {
			argType := fnArgs[i].Type()
			if !argType.AssignableTo(fnArgType) {
				return nil, fmt.Errorf("%d argument %s type mismatch with function %s type", i, argType, fnArgType)
			}
		}
	} else {
		fnArgType := f.Fn.Type().In(f.Fn.Type().NumIn() - 1)
		argType := fnArgs[len(fnArgs)-1].Type()
		if !argType.AssignableTo(fnArgType) {
			return nil, fmt.Errorf("%d argument %s type mismatch with function %s type", len(fnArgs)-1, argType, fnArgType)
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
