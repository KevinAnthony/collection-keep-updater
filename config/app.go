package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/viper"
)

type App struct {
	Series    []types.Series          `json:"series"`
	Libraries []types.LibrarySettings `json:"libraries"`
}

func InitConfig() (App, error) {
	var cfg App

	file, err := os.Open(viper.ConfigFileUsed())
	if err != nil {
		return cfg, err
	}

	defer file.Close()

	bts, err := io.ReadAll(file)
	if err != nil {
		return cfg, err
	}

	if err = json.Unmarshal(bts, &cfg); err != nil {
		return cfg, err
	}

	if err = env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
