// config/config.go
package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aasoft24/golara/wpkg/configs" // 👉 এখানে তোমার module path লাগবে
)

// Init wrapper
func Init(path string) {
	// wpkg/configs এর LoadConfig কল করবে
	configs.LoadConfig(path)

	// 2️⃣ Override with root .env if exists
	if _, err := os.Stat(".env"); err == nil {
		// Load .env manually
		loadRootEnv()
	}
}

// Direct getter
func Get() *configs.Config {
	return configs.GConfig
}

// ------------------- helper -------------------
func loadRootEnv() {
	// App
	if val := os.Getenv("APP_NAME"); val != "" {
		configs.GConfig.App.Name = val
	}
	if val := os.Getenv("APP_ENV"); val != "" {
		configs.GConfig.App.Env = val
	}
	if val := os.Getenv("APP_TIMEZONE"); val != "" {
		configs.GConfig.App.Timezone = val
	}

	// Server
	if val := os.Getenv("SERVER_HOST"); val != "" {
		configs.GConfig.Server.Host = val
	}
	if val := os.Getenv("SERVER_PORT"); val != "" {
		configs.GConfig.Server.Port = atoiSafe(val, configs.GConfig.Server.Port)
	}

	// Redis
	if val := os.Getenv("REDIS_HOST"); val != "" {
		configs.GConfig.Redis.Host = val
	}
	if val := os.Getenv("REDIS_PORT"); val != "" {
		configs.GConfig.Redis.Port = atoiSafe(val, configs.GConfig.Redis.Port)
	}
	if val := os.Getenv("REDIS_PASSWORD"); val != "" {
		configs.GConfig.Redis.Password = val
	}
	if val := os.Getenv("REDIS_DB"); val != "" {
		configs.GConfig.Redis.DB = atoiSafe(val, configs.GConfig.Redis.DB)
	}
	if val := os.Getenv("REDIS_ENABLED"); val != "" {
		configs.GConfig.Redis.Enabled = (val == "true" || val == "1")
	}

	// timezone apply
	if configs.GConfig.App.Timezone != "" {
		if loc, err := time.LoadLocation(configs.GConfig.App.Timezone); err == nil {
			time.Local = loc
		}
	}
}

func atoiSafe(s string, def int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}

// Optional type-safe getters
func GetString(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func GetInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func GetBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		val = strings.ToLower(val)
		return val == "true" || val == "1"
	}
	return defaultVal
}
