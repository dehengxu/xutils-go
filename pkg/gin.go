package xutils

import "encoding/json"

func GetJsonBytes[T any](c *T) *[]byte {
	if bytes, err := json.Marshal(c); err == nil {
		return &bytes
	}
	return nil
}
