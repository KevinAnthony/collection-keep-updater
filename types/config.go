package types

import "github.com/spf13/cobra"

type Config struct {
	Series    []Series          `json:"series"    mapstructure:"series"`
	Libraries []LibrarySettings `json:"libraries" mapstructure:"libraries"`
}

func getFlagOrDefault[T any](cmd *cobra.Command, key string, flagValue, defaultValue T) T {
	flag := cmd.Flag(key)
	if flag == nil {
		return defaultValue
	}

	if flag.Changed {
		return flagValue
	}

	return defaultValue
}
