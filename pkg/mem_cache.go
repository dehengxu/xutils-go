package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	bigcache "github.com/allegro/bigcache/v3"
)

var memCache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))

func SaveMemCache(key string, data interface{}, expiration time.Duration) {
	if bytes, err := json.Marshal(data); err != nil {
		fmt.Printf("Failed to caching, Marshal data, error: %v\n", err)
	} else {
		keyExpir := fmt.Sprintf("%v_expir", key)
		memCache.Set(key, bytes)
		if expiration != 0 {
			memCache.Set(keyExpir, []byte(time.Now().Add(expiration).Format(time.RFC3339)))
		} else {
			memCache.Delete(keyExpir)
		}
	}
}

func LoadMemCache[T any](key string) T {
	var data T
	bytes, err := memCache.Get(key)
	if err != nil {
		return data
	}

	if t, err := memCache.Get(key + "_expir"); err == nil {
		expir, _ := time.Parse(time.RFC3339, string(t))
		if time.Now().After(expir) {
			memCache.Delete(key)
			return data
		}
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		fmt.Printf("Failed to caching, Unmarshal data, error: %v\n", err)
	}
	return data
}
