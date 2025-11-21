package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	// 测试连接
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis连接失败: %v", err)
	}

	fmt.Println("Redis连接成功")
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// SetString 设置字符串值
func SetString(key, value string, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// GetString 获取字符串值
func GetString(key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Get(ctx, key).Result()
}

// Delete 删除键
func Delete(key string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	if RedisClient == nil {
		return false, fmt.Errorf("Redis客户端未初始化")
	}
	result, err := RedisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// 👈 新增：GetTTL 获取键的剩余过期时间
func GetTTL(key string) (time.Duration, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.TTL(ctx, key).Result()
}

// 👈 新增：SetExpire 设置键的过期时间
func SetExpire(key string, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Expire(ctx, key, expiration).Err()
}

// 👈 新增：GetKeys 获取匹配模式的所有键
func GetKeys(pattern string) ([]string, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Keys(ctx, pattern).Result()
}

// 👈 新增：Increment 递增计数器
func Increment(key string) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Incr(ctx, key).Result()
}

// 👈 新增：Decrement 递减计数器
func Decrement(key string) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Decr(ctx, key).Result()
}

// 👈 新增：SetHash 设置哈希字段
func SetHash(key, field, value string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.HSet(ctx, key, field, value).Err()
}

// 👈 新增：GetHash 获取哈希字段
func GetHash(key, field string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.HGet(ctx, key, field).Result()
}

// 👈 新增：GetAllHash 获取哈希的所有字段
func GetAllHash(key string) (map[string]string, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.HGetAll(ctx, key).Result()
}
