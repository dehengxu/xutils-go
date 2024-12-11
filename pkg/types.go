package pkg

import "reflect"

func IsBool(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Bool
}

func IsError(v interface{}) bool {
	_, ok := v.(error)
	return ok
}
