package core

import (
	"gopkg.in/yaml.v3"
	"log"
	"server/config"
	"server/utils"
)

func InitConfig() *config.Config {
	c := &config.Config{}
	yamlConfig, err := utils.LoadYAML()
	if err != nil {
		log.Fatalf("Failed to load yaml config file: %v", err)
	}
	err = yaml.Unmarshal(yamlConfig, c)
	if err != nil {
		log.Fatalf("Failed to unmarshal yaml config file: %v", err)
	}

	return c
}
