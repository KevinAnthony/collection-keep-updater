package config

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configurations",
}

var (
	series  bool
	library bool
)

func init() {
	Cmd.PersistentFlags().BoolVarP(&series, "series", "s", false, "List one or all series configurations")
	Cmd.PersistentFlags().BoolVarP(&library, "library", "l", false, "List one or all library configurations")

	Cmd.AddCommand(listCmd, addCmd, removeCmd)
}
