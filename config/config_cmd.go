package config

import (
	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/spf13/cobra"
)

const (
	seriesFlag  = "series"
	libraryFlag = "library"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configurations",
}

var (
	isSeries  bool
	isLibrary bool
)

func GetCmd() *cobra.Command {
	return configCmd
}

func init() {
	configCmd.PersistentFlags().BoolVarP(&isSeries, seriesFlag, "s", false, "List one or all series configurations")
	configCmd.PersistentFlags().BoolVarP(&isLibrary, libraryFlag, "l", false, "List one or all library configurations")
	configCmd.MarkFlagsOneRequired(seriesFlag, libraryFlag)
	configCmd.MarkFlagsMutuallyExclusive(seriesFlag, libraryFlag)

	configCmd.AddCommand(addCmd, listCmd, removeCmd)

	types.SeriesSetFlags(configCmd)
}
