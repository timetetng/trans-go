// Package config used for load configs
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
}

func InitConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".trans")

	viper.SetDefault("base_url", "https://api.openai.com/v1")
	viper.SetDefault("model", "gpt-3.5-turbo")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			fmt.Println("读取配置文件出错:", err)
		}
	}
}

func SaveConfig(key, value string) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}

func CreateConfigFile() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".trans.yaml")
	return viper.WriteConfigAs(path)
}

func GetConfig() *Config {
	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil
	}
	return &c
}
