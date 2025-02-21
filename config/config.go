package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config struct untuk menyimpan data konfigurasi
type Config struct {
	Mongo struct {
		URI      string `yaml:"uri"`
		Database string `yaml:"database"`
	} `yaml:"mongo"`
}

// Fungsi untuk memuat konfigurasi dari file config.yaml
func LoadConfig() (*Config, error) {
	file, err := os.Open("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("error membuka file config.yaml: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error membaca file config.yaml: %v", err)
	}
	return &config, nil
}
