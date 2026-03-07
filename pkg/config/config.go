package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

var App Config

type Config struct {
	Env      string   `yaml:"env"`
	Api      Api      `yaml:"api"`
	Database Database `yaml:"database"`
}

type Api struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Database struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Keyspace    string `yaml:"keyspace"`
	Consistency string `yaml:"consistency"`
}

func Init(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &App)
}
