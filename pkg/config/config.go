package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SourceURL   string `yaml:"source_url"`
	DBFile      string `yaml:"db_file"`
	IndexFile   string `yaml:"index_file"`
	Parallel    int    `yaml:"parallel"`
	SearchLimit int    `yaml:"search_limit"`
}

func ReadConfig(filename string) (*Config, error) {
	// Считывание конфига из config.yaml
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	return config, nil

}
