package pkg

import (
	"os"
	"reflect"
	"strconv"
)

func GetEnv[T any](key string, defaultValue T) T {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var convertedValue any
	switch reflect.TypeOf(defaultValue).Kind() {
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		convertedValue = intValue
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return defaultValue
		}
		convertedValue = floatValue
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		convertedValue = boolValue
	case reflect.String:
		convertedValue = value
	default:
		// Return the default value for unsupported types
		return defaultValue
	}

	return convertedValue.(T)
}
