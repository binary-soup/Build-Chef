package config

import (
	"encoding/json"
	"errors"
	"os"
)

const CONFIG_FILE = "config.json"

func Load() (Config, error) {
	cfg := Config{}

	file, err := os.Open(CONFIG_FILE)
	if err != nil {
		return cfg, errors.Join(errors.New("error opening file"), err)
	}

	defer file.Close()
	d := json.NewDecoder(file)

	err = d.Decode(&cfg)
	if err != nil {
		return cfg, errors.Join(errors.New("error decoding json"), err)
	}

	return cfg, nil
}

type Config struct {
	Compiler    string   `json:"compiler"`
	SystemPaths []string `json:"system_paths"`
}
