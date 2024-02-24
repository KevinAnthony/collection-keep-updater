package config

import (
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new configuration",
	RunE:  runAdd,
}

func init() {
	types.SeriesSetFlags(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	return nil
}
