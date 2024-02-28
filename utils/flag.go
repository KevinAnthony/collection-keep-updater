package utils

import "github.com/spf13/cobra"

func GetFlagOrDefault[T any](cmd *cobra.Command, key string, flagValue, defaultValue T) T {
	flag := cmd.Flag(key)
	if flag == nil {
		return defaultValue
	}

	if flag.Changed {
		return flagValue
	}

	return defaultValue
}
