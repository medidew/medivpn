package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server_ip   string `yaml:"medivpn_server_ip"`
	Server_port string `yaml:"medivpn_server_port"`
}

func LoadConfig() (*Config, error) {
	home_dir := os.Getenv("HOME")
	if len(home_dir) == 0 {
		return nil, fmt.Errorf("HOME environment variable is empty")
	}
	config_path := filepath.Join(home_dir, ".config", "medivpn", "config.yaml")

	config_file, err := os.Open(config_path)
	if err != nil {
		return nil, fmt.Errorf("error opening config: %v", err)
	}
	defer config_file.Close()

	config_raw, err := io.ReadAll(config_file)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	cfg := new(Config)
	if err = yaml.Unmarshal(config_raw, cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	return cfg, nil
}
