package config

import (
	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Application struct {
	Auth struct {
		Username string `env:"AUTH_USERNAME,required"`
		Password string `env:"AUTH_PASSWORD,required"`
	}
	Path struct {
		UploadedDir  string
		ProcessedDir string
		Cwd          string
	}
}

type ServerConfig struct {
	Constants struct {
		CertFile    string `yaml:"certFile"`
		KeyFile     string `yaml:"keyFile"`
		ServerIP    string `yaml:"serverIP"`
		ServerPort  string `yaml:"serverPort"`
		FileStorage string `yaml:"fileStorage"`
	} `yaml:"constants"`
	ConfigFile string `env:"CONFIG" envDefault:"../../internal/config/resources/defaultConfig.yaml"`
}

func NewConfiguration() (*ServerConfig, *Application, error) {
	cfg := ServerConfig{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, nil, err
	}
	yamlFile, err := ioutil.ReadFile(cfg.ConfigFile)
	if err != nil {
		return nil, nil, err
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, nil, err
	}
	app := Application{}
	err = env.Parse(&app)
	if err != nil {
		return nil, nil, err
	}
	app.Path.UploadedDir = cfg.Constants.FileStorage + "/uploaded-files/"
	app.Path.ProcessedDir = cfg.Constants.FileStorage + "/processed-files/"
	app.Path.Cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &cfg, &app, nil
}
