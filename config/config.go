package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const DefaultConfigPath = "./config.json"

type Config struct {
	Chains []RawChainConfig `json:"chains"`
	Other  Construction     `json:"other,omitempty"`
}

// RawChainConfig is parsed directly from the config file and should be using to construct the core.ChainConfig
type RawChainConfig struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Id           string `json:"id"`       // ChainID
	Endpoint     string `json:"endpoint"` // url for rpc endpoint
	KeystorePath string `json:"keystorePath"`
	Opts         opt    `json:"opts"`
}

type Construction struct {
	MonitorUrl string `json:"monitor_url,omitempty"`
	Etcd       string `json:"etcd,omitempty"`
	Redis      string `json:"redis,omitempty"`
	Env        string `json:"env,omitempty"`
	Db         string `json:"db,omitempty"`
}

type opt struct {
	Mcs                string `json:"mcs,omitempty"`
	StartBlock         string `json:"startBlock,omitempty"`
	Event              string `json:"event,omitempty"`
	BlockConfirmations string `json:"blockConfirmations,omitempty"`
}

func (c *Config) validate() error {
	for _, chain := range c.Chains {
		if chain.Id == "" {
			return fmt.Errorf("required field chain.Id empty for chain %s", chain.Id)
		}
		if chain.Type == "" {
			return fmt.Errorf("required field chain.Type empty for chain %s", chain.Id)
		}
		if chain.Endpoint == "" {
			return fmt.Errorf("required field chain.Endpoint empty for chain %s", chain.Id)
		}
		if chain.Name == "" {
			return fmt.Errorf("required field chain.Name empty for chain %s", chain.Id)
		}
	}
	return nil
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

	err = fig.validate()
	if err != nil {
		return nil, err
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
