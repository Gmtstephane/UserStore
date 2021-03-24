package config

import (
	"log"

	"gopkg.in/yaml.v2"
)

func GenerateConfigFromFile(config []byte) Config {
	c := Config{}
	//yamlFile, err := ioutil.ReadFile(path)
	//if err != nil {
	//	log.Fatalf("yamlFile.Get err   #%v ", err)
	//	}
	err := yaml.Unmarshal(config, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	ValidateConfig(c)
	return c
}

func ValidateConfig(config Config) {
	if config.Database.Name == "" {
		log.Fatalf("No Database provided")
	} else if config.Database.Uri == "" {
		log.Fatalf("No Uri provided")
	} else if config.Database.Port == 0 {
		log.Fatalf("No Port provided")
	}
}
