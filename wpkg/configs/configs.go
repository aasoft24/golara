// wpkg/configs/configs.go
package configs

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Name     string `yaml:"name"`
	Env      string `yaml:"env"`
	Timezone string `yaml:"timezone"`
}

type Config struct {
	App      AppConfig
	Database struct {
		Default     string
		Connections map[string]map[string]string
	}
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
		Enabled  bool   `yaml:"enabled"`
	} `yaml:"redis"`
}

var GConfig *Config

func LoadConfig(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Config file read error: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalf("YAML parse error: %v", err)
	}

	GConfig = &cfg

	// 🔥 Global timezone set
	if GConfig.App.Timezone != "" {
		loc, err := time.LoadLocation(GConfig.App.Timezone)
		if err != nil {
			log.Fatalf("Invalid timezone: %v", err)
		}
		time.Local = loc
		log.Printf("🌍 Timezone set to %s", GConfig.App.Timezone)
	}
}

func Init() {
	LoadConfig("config.yaml")
}
