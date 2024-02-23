package config

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a configuration",
}
