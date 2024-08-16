package config

import (
	"os"
	"path/filepath"

	"gopkg.in/gcfg.v1"
)

type Config struct {
	FaseRule map[string]*struct {
		Price             float64
		Distance          float64
		DistanceThreshold float64
	}
}

var readFile = gcfg.ReadFileInto

func InitMainConfig() (*Config, error) {
	cfg := Config{}
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return &cfg, err
	}

	err = readFile(&cfg, path+"/config-main.ini")
	if err != nil {
		return &cfg, err
	}

	return &cfg, err
}
