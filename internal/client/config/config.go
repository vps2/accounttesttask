package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Client struct {
		Addr    string `yaml:"addr"`
		Readers int    `yaml:"readers"`
		Writers int    `yaml:"writers"`
		Keys    []int  `yaml:"keys"`
	}	`yaml:"client"`
}

func New(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	cfg := Config{}
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
