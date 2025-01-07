package pkg

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestRedis(t *testing.T) {
	// 初始化 Redis 客户端
	InitRedis(nil)

	ctx := &RedisContext{
		Conf: &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	}

	// 测试保存数据
	user := User{
		Name: "John Doe",
		Age:  30,
	}

	err := SaveToRedis(ctx, "user:123", user, 10*time.Minute)
	if err != nil {
		fmt.Printf("Error saving to Redis:%v, %v \n", err, user)
		return
	}
	fmt.Printf("Error saving to Redis: %v \n", user)

	// 测试加载数据
	var loadedUser User
	loadedUser, err = LoadFromRedis[User](ctx, "user:123")
	if err != nil {
		fmt.Println("Error loading from Redis:", err)
		return
	}
	fmt.Println("Data loaded from Redis:", loadedUser)
}
