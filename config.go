package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Urls []string
}

func ParseConfig() Config {
	var config Config
	config_path := "config.yml"
	data, _ := ioutil.ReadFile(config_path)

	err := yaml.Unmarshal(data, &config)

	if err != nil {
		panic("Error reading config.yml")
	}

	return config
}
