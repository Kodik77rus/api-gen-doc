package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port         string `mapstructure:"PORT"`
	ReadTimeout  string `mapstructure:"READ_TIMEOUT"`
	WriteTimeout string `mapstructure:"WRITE_READ_TIMEOUT"`
}

type TemplateBuilder struct {
	WordFolder string `mapstructure:"WORD_FOLDER"`
}

type AppConfig struct {
	Server          ServerConfig    `mapstructure:",squash"`
	TemplateBuilder TemplateBuilder `mapstructure:",squash"`
}

type errNotFound struct {
	file string
}

func (e errNotFound) Error() string {
	return fmt.Sprintf("no templates in the folder %s", e.file)
}

// App config constructor
func LoadConfig() (*AppConfig, error) {
	var result map[string]interface{}
	var config AppConfig

	viper.SetConfigFile("../.env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&result); err != nil {
		return nil, err
	}

	if err := mapstructure.Decode(result, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
