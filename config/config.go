/*
	Copyright © 2022 Tom Lister tom@tomlister.net
*/
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Key      string `json:"key"`
}

func (c *Config) Save() error {
	// make a config file in the users platform specific config dir

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	blikDir := path.Join(cfgDir, "/blik-cli")
	err = os.Mkdir(blikDir, 0755)
	if !errors.Is(err, os.ErrExist) && err != nil {
		return err
	}

	// marshal the config struct into json
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(blikDir, "config.json"), b, 0755)
	if err != nil {
		return err
	}

	return nil
}

func NewConfig(endpoint string, key string) *Config {
	return &Config{
		Endpoint: endpoint,
		Key:      key,
	}
}

func ReadConfig() (*Config, error) {
	var cfg Config

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return &cfg, err
	}

	cfgPath := path.Join(cfgDir, "/blik-cli/", "config.json")
	f, err := os.Open(cfgPath)
	if err != nil {
		return &cfg, err
	}

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
