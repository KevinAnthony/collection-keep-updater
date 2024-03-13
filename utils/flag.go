package utils

import (
	"github.com/kevinanthony/collection-keep-updater/types"
)

func GetFlagOrDefault[T any](cmd types.ICommand, key string, flagValue, defaultValue T) T {
	flag := cmd.Flag(key)
	if flag == nil {
		return defaultValue
	}

	if flag.Changed {
		return flagValue
	}

	return defaultValue
}
