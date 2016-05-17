package main

import (
	"log"

	"gopkg.in/yaml.v2"
)

type Command struct {
	Image   string `image`
	Command string `command`
}

func Parse(text string) map[string]Command {
	bytes := []byte(text)
	config := make(map[string]Command)

	err := yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return config
}
