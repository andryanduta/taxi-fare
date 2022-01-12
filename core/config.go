package config

import (
	"os"
	"path/filepath"

	"gopkg.in/gcfg.v1"
)

type Config struct {
	FaseRules map[string]*struct {
		Price             float64
		Distance          int
		DistanceThreshold int
	}
}

func InitMainConfig() (*Config, error) {
	cfg := Config{}
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return &cfg, err
	}

	err = gcfg.ReadFileInto(&cfg, path+"/config-main.ini")
	if err != nil {
		return &cfg, err
	}

	return &cfg, err
}
