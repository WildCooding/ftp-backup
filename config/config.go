package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Username             string   `yaml:"username"`
	Password             string   `yaml:"password"`
	Host                 string   `yaml:"host"`
	Port                 string   `yaml:"port"`
	CrontabInterval      string   `yaml:"crontab_interval"`
	Directorys           []string `yaml:"directorys"`
	DestinationDirectory string   `yaml:"destination_directory"`
}

func LoadConfig() (*Config, error) {

	username := os.Getenv("username")
	password := os.Getenv("password")
	host := os.Getenv("host")
	port := os.Getenv("port")
	contabInterval := os.Getenv("crontab_interval")
	directorys := os.Getenv("directorys")
	destinationDirectory := os.Getenv("destination_directory")

	directorys_array := strings.Split(directorys, ",")

	if username != "" && password != "" && host != "" && port != "" && contabInterval != "" && destinationDirectory != "" && directorys != "" && len(directorys_array) > 0 {
		return &Config{
			Username:             username,
			Password:             password,
			Host:                 host,
			Port:                 port,
			CrontabInterval:      contabInterval,
			Directorys:           directorys_array,
			DestinationDirectory: destinationDirectory,
		}, nil
	}

	config := &Config{}
	yamlData, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(yamlData, &config)

	return config, nil
}
