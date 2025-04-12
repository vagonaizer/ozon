package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// PSConfig - (Product-Service-Config)
	PSConfig struct {
		ProductServiceURL   string `yaml:"product_service_url"`
		ProductServiceToken string `yaml:"product_service_token"`
	} `yaml:"ps_config"`
	ServerConfig struct {
		Port   string `yaml:"port"`
		BindIP string `yaml:"bind_ip"`
	} `yaml:"server_config"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() (*Config, error) {
	once.Do(func() {
		data, err := os.ReadFile("config.yml")
		if err != nil {
			return
		}

		var cfg Config
		if err = yaml.Unmarshal(data, &cfg); err != nil {
			return
		}
		instance = &cfg
	})
	return instance, nil
}
