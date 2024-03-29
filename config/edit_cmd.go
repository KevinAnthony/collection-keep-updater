package config

import (
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit an existing configuration",
	RunE:  types.CmdRunE(runEdit),
}

func init() {
	editCmd.PersistentFlags().BoolVarP(&try, "test-config", "t", false, "test the configuration by calling source and outputting result.")
	editCmd.PersistentFlags().BoolVarP(&write, "write-config", "w", false, "save the configuration.")

	editCmd.MarkFlagsOneRequired("test-config", "write-config")
	editCmd.MarkFlagsMutuallyExclusive("test-config", "write-config")

	seriesSetFlags(editCmd)
}

func runEdit(cmd types.ICommand, args []string) error {
	cfg, err := ctxu.GetConfig(cmd)
	if err != nil {
		return err
	}

	switch {
	case isSeries:
		s, err := editSeries(cmd, cfg)
		if err != nil {
			return err
		}

		if try {
			source, err := ctxu.GetSource(cmd, s.Source)
			if err != nil {
				return err
			}

			books, err := source.GetISBNs(cmd.Context(), *s)
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

			cfg.Series = append(cfg.Series, *s)
			viper.Set("series", cfg.Series)

			return viper.WriteConfig()
		}
	}

	return nil
}
