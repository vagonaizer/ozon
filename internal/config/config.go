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
	var cfgErr error
	once.Do(func() {
		file, err := os.Open("config.yml")
		if err != nil {
			cfgErr = err
			return
		}
		defer file.Close()

		var cfg Config
		if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
			cfgErr = err
			return
		}
		instance = &cfg
	})

	if cfgErr != nil {
		return nil, cfgErr
	}
	return instance, nil
}
