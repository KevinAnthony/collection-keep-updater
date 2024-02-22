package config

import (
	"github.com/kevinanthony/collection-keep-updater/types"
)

type App struct {
	Series    []types.Series          `json:"series"    mapstructure:"series"`
	Libraries []types.LibrarySettings `json:"libraries" mapstructure:"libraries"`
}
