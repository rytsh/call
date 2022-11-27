package call

import (
	"reflect"
	"strings"
	"sync"
)

type Option interface {
	GetDelimeter() string
	AddOption(name string, fn func([]reflect.Value, ...string) ([]reflect.Value, error)) Option
	GetOption(name string) (func([]reflect.Value, ...string) ([]reflect.Value, error), bool)
	VisitOptions(arg string, v any) ([]reflect.Value, error)
}

type OptionFunc struct {
	Name string
	Fn   func([]reflect.Value, ...string) ([]reflect.Value, error)
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
func NewReg(optionFuncs ...OptionFunc) *Reg {
	option := NewOptions().
		AddOption("index", OptionGetIndex).
		AddOption("...", OptionVariadic)

	for _, o := range optionFuncs {
		option.AddOption(o.Name, o.Fn)
	}

	return &Reg{
		fn:     make(map[string]Func),
		args:   make(map[string]any),
		Option: option,
	}
}

// AddArgument adds argument to registry with name.
//
// If name includes delimeter, it will not add options.
func (r *Reg) AddArgument(name string, v any) *Reg {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// trim options
	name = strings.SplitN(name, r.GetDelimeter(), 2)[0]

	r.args[name] = v

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
		panic("fn argument is not a function")
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

// DeleteFunction removes function with name.
func (r *Reg) DeleteFunction(name string) *Reg {
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
