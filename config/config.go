package config

import (
	"log"
	"os"
	"path/filepath"
	"queue/database"
	"queue/schema"
	"queue/usecase"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	DB      database.Config
	Usecase usecase.Config
	Schema  schema.Config
}

func New(path string) *Conf {

	var conf Conf
	var model schema.Model

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // or viper.SetConfigType("YAML")
	viper.AddConfigPath(path)     // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %+v \n", err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	err = readYAMLFile(path, "config.model.yaml", &model)
	if err != nil {
		log.Fatalf("Unable to read yaml file, %v", err)
	}
	conf.Schema.Model = model

	return &conf
}

func readYAMLFile(path, filename string, data interface{}) error {
	fullPath := filepath.Join(path, filename)
	yamlFile, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		return err
	}

	return nil
}
