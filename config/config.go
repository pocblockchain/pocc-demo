package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	DefaultGasLimit     = 200000
	DefaultMiniGasPrice = "10000000000poc"
)

// Config holds all configurations.
type Config struct {
	RPCURL      string `yaml:"rpc_url"`
	ChainID     string `yaml:"chain_id"`
	Gas         uint64 `yaml:"gas"`
	MinGasPrice string `yaml:"min_gasprice"`
}

// GetConfig parses the config file into a Config instance.
func GetConfig(path string) (*Config, error) {
	var config = new(Config)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	if config.RPCURL == "" || config.ChainID == "" ||
		config.MinGasPrice == "" {
		return nil, fmt.Errorf("configuration is null")
	}

	if config.Gas == 0 {
		config.Gas = DefaultGasLimit
	}

	return config, nil
}
