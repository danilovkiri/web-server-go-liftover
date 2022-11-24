// Package config provides ENV and YAML parsing functionality.
package config

import (
	"github.com/caarlos0/env/v6"
)

// Application defines a set of app parameters and their parsing from ENV.
type Application struct {
	Auth struct {
		Username string `env:"AUTH_USERNAME,required"`
		Password string `env:"AUTH_PASSWORD,required"`
	}
	Path struct {
		UploadedDir  string
		ProcessedDir string
		WD           string `env:"HOME,required"`
	}
}

// ServerConfig defines a set of server parameters and their parsing from YAML.
type ServerConfig struct {
	Constants struct {
		CertFile    string `env:"CERT,required"`
		KeyFile     string `env:"KEY,required"`
		FileStorage string `env:"STORAGE" envDefault:"./file-storage"`
		ServerHost  string `env:"HOST" envDefault:"0.0.0.0"`
		ServerPort  string `env:"PORT" envDefault:"8080"`
	}
}

// NewConfiguration parses YAML and ENV filling the configuration object.
func NewConfiguration() (*ServerConfig, *Application, error) {
	cfg := ServerConfig{}
	err := env.Parse(&cfg)
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
	if err != nil {
		return nil, nil, err
	}
	return &cfg, &app, nil
}
