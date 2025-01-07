package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// var rdb *redis.Client
var ctx = context.Background()
var rdbs = make(map[string]*redis.Client)

type RedisContext struct {
	Conf *redis.Options
}

func (c *RedisContext) GetKey() string {
	return fmt.Sprintf("%v-%d", c.Conf.Addr, c.Conf.DB)
}

// 初始化 Redis 客户端
func InitRedis(conf *redis.Options) {
	key := ""
	if conf == nil {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379", // Redis 服务器地址
			Password: "",               // Redis 密码，如果没有可以留空
			DB:       0,                // 默认数据库 0
		})
		key = fmt.Sprintf("localhost:6379-%d", 0)
		rdbs[key] = rdb
	} else {
		key = fmt.Sprintf("%v-%d", conf.Addr, conf.DB)
		rdb := redis.NewClient(conf)
		rdbs[key] = rdb
	}
}

// LoadFromRedis 从 Redis 加载数据
func LoadFromRedis[T any](c *RedisContext, key string) (T, error) {
	var result T
	rdb, exists := rdbs[c.GetKey()]
	if !exists {
		return result, fmt.Errorf("redis client not found")
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Redis 中没有数据
			return result, fmt.Errorf("key not found")
		}
		return result, err
	}
	fmt.Printf("load val:%v\n", val)
	// 将缓存的数据反序列化
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal data: %v", err)
	}
	return result, nil
}

// SaveToRedis 保存数据到 Redis
func SaveToRedis(c *RedisContext, key string, data interface{}, expir time.Duration) error {
	rdb, exists := rdbs[c.GetKey()]
	if !exists {
		return fmt.Errorf("redis client not found")
	}
	// 序列化数据
	val, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	fmt.Printf("save val:%v\n", val)

	// 将数据保存到 Redis，设置过期时间
	err = rdb.Set(ctx, key, val, expir).Err()
	if err != nil {
		return fmt.Errorf("failed to save data to redis: %v", err)
	}
	return nil
}
