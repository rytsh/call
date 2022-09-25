package call

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Options struct {
	option map[string]func([]reflect.Value, ...string) ([]reflect.Value, error)
	mutex  sync.RWMutex
}

var _ Option = (*Options)(nil)

func NewOptions() Option {
	return &Options{
		option: make(map[string]func([]reflect.Value, ...string) ([]reflect.Value, error)),
	}
}

// GetDelimeter returns delimeter for options.
func (o *Options) GetDelimeter() string {
	return ":"
}

// ParseOptions parses options from string with delimeter.
func (o *Options) ParseOption(name string) []string {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	sName := strings.SplitN(name, ":", 2)

	if len(sName) == 1 {
		return nil
	}

	return strings.Split(sName[1], ";")
}

func (o *Options) AddOption(name string, fn func([]reflect.Value, ...string) ([]reflect.Value, error)) Option {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.option[name] = fn

	return o
}

func (o *Options) GetOption(name string) (func([]reflect.Value, ...string) ([]reflect.Value, error), bool) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	fn, ok := o.option[name]

	return fn, ok
}

func (o *Options) VisitOptions(arg string, v any) (error, []reflect.Value) {
	var err error

	vValue := []reflect.Value{reflect.ValueOf(v)}

	options := o.ParseOption(arg)

	for _, option := range options {
		opt := strings.SplitN(option, "=", 2)
		optName := opt[0]

		var optVariables []string

		if len(opt) > 1 {
			optVariables = strings.Split(opt[1], ",")
		}

		fn, ok := o.GetOption(optName)
		if !ok {
			return fmt.Errorf("option %s not found", optName), nil
		}

		vValue, err = fn(vValue, optVariables...)
		if err != nil {
			return fmt.Errorf("%s; %w", optName, err), nil
		}

		if vValue == nil {
			break
		}
	}

	return nil, vValue
}
