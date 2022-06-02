package config

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port         string `mapstructure:"PORT"`
	ReadTimeout  string `mapstructure:"READ_TIMEOUT"`
	WriteTimeout string `mapstructure:"WRITE_READ_TIMEOUT"`
}

type AppConfig struct {
	Server         ServerConfig `mapstructure:",squash"`
	TemplateFolder string       `mapstructure:"TEMPLATE_FOLDER"`
}

type errNotFound struct {
	file string
}

func (e errNotFound) Error() string {
	return fmt.Sprintf("no templates in the folder %s", e.file)
}

// App config constructor
func NewConfig() (*AppConfig, error) {
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

	if err := checkTemplateFolder(config.TemplateFolder); err != nil {
		return nil, err
	}

	return &config, nil
}

func checkTemplateFolder(path string) error {
	templatesFolder, err := os.Open(filepath.ToSlash(path))
	if err != nil {
		return err
	}
	defer templatesFolder.Close()

	templates, err := templatesFolder.ReadDir(0)
	if err != nil {
		return err
	}

	if len(templates) == 0 {
		return errNotFound{file: path}
	}

	for _, tempale := range templates {
		fmt.Println(tempale.Name(), "- tempalte loaded")
	}

	return nil
}
