package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Timeout time.Duration
	Urls    []string
}

func ParseConfig() Config {
	var config Config
	config_path := "config.yml"
	data, _ := ioutil.ReadFile(config_path)

	err := yaml.Unmarshal(data, &config)

	if err != nil {
		panic("Error reading config.yml")
	}

	if len(config.Urls) == 0 {
		fmt.Println("You have no urls defined. Double check config.yml")
		os.Exit(2)
	}

	return config
}
