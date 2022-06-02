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

type TemplateBuilder struct {
	TemplateFolder string `mapstructure:"TEMPLATE_FOLDER"`
	WordFolder     string `mapstructure:"WORD_FOLDER"`
	PdfFolder      string `mapstructure:"PDF_FOLDER"`

	Tempaltes []string //files name
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

	if err := config.setTemplates(config.TemplateBuilder.TemplateFolder); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *AppConfig) setTemplates(path string) error {
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

	for _, template := range templates {
		c.SetTemplate(template.Name())
		fmt.Println(template.Name(), "- tempalte loaded")
	}

	return nil
}

func (c *AppConfig) SetTemplate(t string) {
	c.TemplateBuilder.Tempaltes = append(c.TemplateBuilder.Tempaltes, t)
}
