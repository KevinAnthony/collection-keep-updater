package types

import "github.com/spf13/cobra"

type Config struct {
	Series    []Series          `json:"series"    mapstructure:"series"`
	Libraries []LibrarySettings `json:"libraries" mapstructure:"libraries"`
}

func getFlagOrNil[T any](cmd *cobra.Command, key string, value T) (out T) {
	flag := cmd.Flag(key)
	if flag == nil {
		return out
	}

	if flag.Changed {
		return value
	}

	return out
}
