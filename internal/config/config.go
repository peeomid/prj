package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Folders    []string `json:"folders"`
	CutoffDays int      `json:"cutoff_days"`
}

func DefaultConfig() *Config {
	return &Config{
		Folders:    []string{},
		CutoffDays: 240,
	}
}

func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".prj")
}

func Path() string {
	return filepath.Join(Dir(), "config.json")
}

func Load() (*Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.CutoffDays == 0 {
		cfg.CutoffDays = 240
	}
	return cfg, nil
}

func Save(cfg *Config) error {
	if err := os.MkdirAll(Dir(), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0644)
}

func (c *Config) AddFolder(folder string) bool {
	for _, f := range c.Folders {
		if f == folder {
			return false
		}
	}
	c.Folders = append(c.Folders, folder)
	return true
}

func (c *Config) RemoveFolder(folder string) bool {
	for i, f := range c.Folders {
		if f == folder {
			c.Folders = append(c.Folders[:i], c.Folders[i+1:]...)
			return true
		}
	}
	return false
}
