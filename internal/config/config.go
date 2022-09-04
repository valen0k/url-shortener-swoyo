package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
	Storage struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"storage"`
}

func NewConfig(configFile string, d bool) (*Config, error) {
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if d {
		if err = json.Unmarshal(file, config); err != nil {
			return nil, err
		}
	} else {
		server := struct {
			Server struct {
				Host string `json:"host"`
				Port string `json:"port"`
			} `json:"server"`
		}{}
		if err = json.Unmarshal(file, &server); err != nil {
			return nil, err
		}
		config.Server = server.Server
	}

	return config, nil
}
