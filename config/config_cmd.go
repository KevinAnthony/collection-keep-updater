package config

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configurations",
}

func init() {
	Cmd.AddCommand(listCmd, addCmd, removeCmd)
}
