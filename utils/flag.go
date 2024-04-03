package utils

import (
	"github.com/kevinanthony/collection-keep-updater/types"
)

func GetFlagString(cmd types.ICommand, key string) string {
	str, err := cmd.PersistentFlags().GetString(key)
	if err != nil {
		// TODO log to console
		return ""
	}

	return str
}

func GetFlagStringPtr(cmd types.ICommand, key string) *string {
	str, err := cmd.PersistentFlags().GetString(key)
	if err != nil {
		// TODO log to console
		return nil
	}

	return &str
}

func GetFlagStringSlice(cmd types.ICommand, key string) []string {
	str, err := cmd.PersistentFlags().GetStringSlice(key)
	if err != nil {
		// TODO log to console
		return nil
	}

	return str
}

func GetFlagInt(cmd types.ICommand, key string) int {
	i, err := cmd.PersistentFlags().GetInt(key)
	if err != nil {
		// TODO log to console
		return 0
	}

	return i
}

func GetFlagIntPtr(cmd types.ICommand, key string) *int {
	i := GetFlagInt(cmd, key)

	return &i
}

func GetFlagIntSlice(cmd types.ICommand, key string) []int {
	i, err := cmd.PersistentFlags().GetIntSlice(key)
	if err != nil {
		// TODO log to console
		return nil
	}

	return i
}

func GetFlagBool(cmd types.ICommand, key string) bool {
	b, err := cmd.PersistentFlags().GetBool(key)
	if err != nil {
		// TODO log to console
		return false
	}

	return b
}
