package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type YamlConfig struct {
	Constants struct {
		CertFile   string `yaml:"certFile"`
		KeyFile    string `yaml:"keyFile"`
		ServerIP   string `yaml:"serverIP"`
		ServerPort string `yaml:"serverPort"`
	} `yaml:"constants"`
}

func ParseConfig(config string) YamlConfig {
	yamlFile, err1 := ioutil.ReadFile(config)
	if err1 != nil {
		log.Fatal("Error opening YAML configuration:", err1)
	}
	var yamlConfig YamlConfig
	err2 := yaml.Unmarshal(yamlFile, &yamlConfig)
	if err2 != nil {
		log.Fatal("Error parsing YAML configuration:", err2)
	}
	return yamlConfig
}
