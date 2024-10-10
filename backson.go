package backson

import (
	"fmt"
	"github.com/buger/jsonparser"
	"reflect"
)

type KindError reflect.Kind

func (ke KindError) Error() string {
	return fmt.Sprintf("values %s is not supported", reflect.Kind(ke).String())
}

type parseFunc func(value []byte, dataType jsonparser.ValueType, offset int, err error)

func ParseArray[T any](data []byte, ch chan<- T, keys ...string) error {
	defer close(ch)
	parser, err := parseItem[T](ch)
	if err != nil {
		return err
	}

	_, err = jsonparser.ArrayEach(data, parser, keys...)

	return err
}

func parseValue[T any](val *T) (parseFunc, error) {
	v := reflect.ValueOf(val).Elem()
	if !v.CanSet() {
		return nil, fmt.Errorf("value if type %T is not settible", val)
	}

	k := v.Kind()
	switch k {
	case reflect.Bool:
		return getWrapper(jsonparser.ParseBoolean, jsonparser.Boolean, v.SetBool)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return getWrapper(jsonparser.ParseInt, jsonparser.Number, v.SetInt)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return getWrapper(func(bytes []byte) (uint64, error) {
			i, err := jsonparser.ParseInt(bytes)
			if err != nil {
				return 0, err
			}
			return uint64(i), nil
		}, jsonparser.Number, v.SetUint)
	case reflect.Float32, reflect.Float64:
		return getWrapper(jsonparser.ParseFloat, jsonparser.Number, v.SetFloat)
	case reflect.String:
		return getWrapper(jsonparser.ParseString, jsonparser.String, v.SetString)
	case reflect.Pointer, reflect.Uintptr:
		return nil, KindError(k)
	default:
		return nil, KindError(k)
	}
}

func getWrapper[T any](get func([]byte) (T, error), expected jsonparser.ValueType, set func(T)) (parseFunc, error) {
	var v T
	return func(data []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}
		if dataType != expected {
			err = fmt.Errorf("expected type %s got %s", expected, dataType)
			return
		}
		v, err = get(data)
		if err != nil {
			return
		}
		set(v)
	}, nil
}

func parseItem[T any](ch chan<- T) (parseFunc, error) {
	var item T
	//zeroPtr := new(T)
	k := reflect.ValueOf(item).Kind()
	switch k {
	case reflect.Interface:
	case reflect.Struct:

	default:
		f, err := parseValue[T](&item)
		if err != nil {
			return nil, err
		}
		return func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			f(value, dataType, offset, err)
			ch <- item
		}, nil
	}
	return nil, nil
}
