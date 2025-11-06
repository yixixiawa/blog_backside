package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// AppConfig 应用配置结构
type AppConfig struct {
	Redis RedisConfig `yaml:"redis"`
	Email EmailConfig `yaml:"email"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var GlobalConfig *AppConfig

// LoadConfig 加载配置文件
func LoadConfig() (*AppConfig, error) {
	data, err := os.ReadFile("yaml/application.yaml")
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	GlobalConfig = &config
	return &config, nil
}
