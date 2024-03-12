package config

import (
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/kevinanthony/collection-keep-updater/ctxu"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	try   bool
	write bool
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new configuration",
	RunE:  types.CmdRunE(runAdd),
}

func init() {
	addCmd.PersistentFlags().BoolVarP(&try, "test-config", "t", false, "test the configuration by calling source and outputting result.")
	addCmd.PersistentFlags().BoolVarP(&write, "write-config", "w", false, "save the configuration.")

	addCmd.MarkFlagsOneRequired("test-config", "write-config")
	addCmd.MarkFlagsMutuallyExclusive("test-config", "write-config")

	seriesSetFlags(addCmd)
}

func runAdd(cmd types.ICommand, args []string) error {
	switch {
	case isSeries:
		s, err := newSeriesConfig(cmd)
		if err != nil {
			return err
		}

		if try {
			sources, err := ctxu.GetSources(cmd)
			if err != nil {
				return err
			}

			source, found := sources[s.Source]
			if !found {
				return fmt.Errorf("source type %s not found in source map", s.Source)
			}

			books, err := source.GetISBNs(cmd.Context(), s)
			if err != nil {
				return err
			}

			books.Print(cmd)
		}

		if write {
			cfg, err := ctxu.GetConfig(cmd)
			if err != nil {
				return err
			}

			cfg.Series = append(cfg.Series, s)
			viper.Set("series", cfg.Series)

			return viper.WriteConfig()
		}
	}

	return nil
}
