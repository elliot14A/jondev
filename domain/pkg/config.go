package pkg

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Hash struct {
		FilePath string `mapstructure:"file_path"`
		Key      string `mapstructure:"key"`
	} `mapstructure:"hash"`
	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	v.SetEnvPrefix("JONDEV")
	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found. Using environment variables only")
		} else {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %v", err)
	}

	return &config, nil
}

func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
