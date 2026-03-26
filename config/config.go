package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	BaseDN         string `json:"base_dn"`
	UserSearchBase string `json:"user_search_base"`
	UsernameAttr   string `json:"username_attr"`
	UseTLS         bool   `json:"use_tls"`
	UseStartTLS    bool   `json:"use_start_tls"`
}

func dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cldap"), nil
}

func Load() (*Config, error) {
	d, err := dir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(d, "config.json"))
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func Save(c *Config) error {
	d, err := dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(d, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d, "config.json"), data, 0600)
}
