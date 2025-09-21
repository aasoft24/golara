// pkg/cache/redis_cache.go
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aasoft24/golara/wpkg/configs"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var RDB *redis.Client
var Enabled bool

// Init Redis client from YAML config
// InitRedis initializes Redis client from config
func InitRedis() {
	redisCfg := configs.GConfig.Redis

	Enabled = redisCfg.Enabled

	if !Enabled {
		log.Println("Redis is disabled in config")
		return
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
		// Optional: adjust pool settings
		PoolSize:     10,
		MinIdleConns: 2,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	if err := RDB.Ping(ctx).Err(); err != nil {
		Enabled = false
		RDB = nil
		log.Printf("Redis connection failed: %v. Caching disabled.\n", err)
		return
	}

	log.Println("Redis connected successfully")
}

// Set value with TTL
func SetRedis(key string, value interface{}, ttl time.Duration) error {

	if !Enabled || RDB == nil {
		return fmt.Errorf("redis is not enabled")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RDB.Set(ctx, key, data, ttl).Err()
}

// Get value into struct/dest
func GetRedis(key string, dest interface{}) error {
	if !Enabled || RDB == nil {
		return fmt.Errorf("redis is not enabled")
	}

	val, err := RDB.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Delete key
func DelRedis(key string) error {
	if !Enabled || RDB == nil {
		return fmt.Errorf("redis is not enabled")
	}
	return RDB.Del(ctx, key).Err()
}

// Exists
func ExistsRedis(key string) (bool, error) {
	if !Enabled || RDB == nil {
		return false, fmt.Errorf("redis is not enabled")
	}
	n, err := RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Increment
func IncrRedis(key string) (int64, error) {
	if !Enabled || RDB == nil {
		return 0, fmt.Errorf("redis is not enabled")
	}
	return RDB.Incr(ctx, key).Result()
}

// Decrement
func DecrRedis(key string) (int64, error) {
	if !Enabled || RDB == nil {
		return 0, fmt.Errorf("redis is not enabled")
	}
	return RDB.Decr(ctx, key).Result()
}

// Flush all
func FlushAllRedis() error {
	if !Enabled || RDB == nil {
		return fmt.Errorf("redis is not enabled")
	}
	return RDB.FlushAll(ctx).Err()
}

// DebugRedis prints info about a Redis key
func DebugRedis(key string) {
	if !Enabled || RDB == nil {
		fmt.Printf("Redis is not enabled or client is nil\n")
		return
	}

	exists, err := ExistsRedis(key)
	if err != nil {
		fmt.Printf("Redis error checking existence: %v\n", err)
		return
	}

	if !exists {
		fmt.Printf("Key '%s' does NOT exist in Redis.\n", key)
		return
	}

	ttl, err := RDB.TTL(ctx, key).Result()
	if err != nil {
		fmt.Printf("Redis error getting TTL: %v\n", err)
		return
	}

	val, err := RDB.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("Redis error getting value: %v\n", err)
		return
	}

	fmt.Printf("Key: %s\nExists: %t\nTTL: %v\nValue: %s\n", key, exists, ttl, val)
}

// Debug multiple keys by pattern
func DebugRedisPattern(pattern string) {
	if !Enabled || RDB == nil {
		fmt.Printf("Redis is not enabled or client is nil\n")
		return
	}
	var cursor uint64
	for {
		keys, nextCursor, err := RDB.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			fmt.Printf("Redis scan error: %v\n", err)
			return
		}

		for _, key := range keys {
			DebugRedis(key)
			fmt.Println("----")
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}
