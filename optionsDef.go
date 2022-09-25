package call

import (
	"fmt"
	"reflect"
	"strconv"
)

func OptionGetIndex(v []reflect.Value, args ...string) ([]reflect.Value, error) {
	if len(v) == 0 {
		return nil, fmt.Errorf("no value")
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("index is empty")
	}

	vValue := v[0]

	switch vValue.Kind() {
	case reflect.Slice, reflect.Array:
		var retV []reflect.Value
		for _, arg := range args {
			index, err := strconv.Atoi(arg)
			if err != nil {
				return nil, fmt.Errorf("index is not a number; %w", err)
			}

			if index < 0 || index >= vValue.Len() {
				return nil, fmt.Errorf("index out of range")
			}

			retV = append(retV, vValue.Index(index))
		}

		return retV, nil
	case reflect.Map:
		var retV []reflect.Value
		for _, arg := range args {
			retV = append(retV, vValue.MapIndex(reflect.ValueOf(arg)))
		}

		return retV, nil
	default:
		return nil, fmt.Errorf("not related type for index")
	}
}

func OptionVariadic(v []reflect.Value, _ ...string) ([]reflect.Value, error) {
	if len(v) > 1 {
		return v, nil
	}

	// check v[0] slice or array or map
	if v[0].Kind() != reflect.Slice && v[0].Kind() != reflect.Array && v[0].Kind() != reflect.Map {
		return nil, fmt.Errorf("not slice or array")
	}

	// get slice
	ret := make([]reflect.Value, v[0].Len())

	if v[0].Kind() == reflect.Map {
		for i, key := range v[0].MapKeys() {
			ret[i] = v[0].MapIndex(key)
		}
	} else {
		for i := 0; i < v[0].Len(); i++ {
			ret[i] = v[0].Index(i)
		}
	}

	return ret, nil
}
