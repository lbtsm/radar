package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultConfigPath = "./config.json"
)

type Config struct {
	Listen string
	Dsn    string
}

func Local(cfgFile string) (*Config, error) {
	var fig Config
	path := DefaultConfigPath
	if cfgFile != "" {
		path = cfgFile
	}

	err := loadConfig(path, &fig)
	if err != nil {
		return &fig, err
	}

	return &fig, nil
}

func loadConfig(file string, config *Config) error {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}

	if ext == ".json" {
		if err = json.NewDecoder(f).Decode(&config); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}

	return nil
}
