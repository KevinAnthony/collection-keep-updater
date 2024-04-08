package config

import (
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit an existing configuration",
	RunE:  types.CmdRunE(runEdit),
}

func init() {
	editCmd.PersistentFlags().BoolP(testF, "t", false, "test the configuration by calling source and outputting result.")
	editCmd.PersistentFlags().BoolP(writeF, "w", false, "save the configuration.")

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
	case utils.GetFlagBool(cmd, seriesFlag):
		s, err := editSeries(cmd, cfg)
		if err != nil {
			return err
		}

		if utils.GetFlagBool(cmd, testF) {
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

		if utils.GetFlagBool(cmd, writeF) {
			cfg, err := ctxu.GetConfig(cmd)
			if err != nil {
				return err
			}

			cfg.Series = append(cfg.Series, *s)
			icfg := ctxu.GetConfigReader(cmd)

			icfg.Set("series", cfg.Series)

			return icfg.WriteConfig()
		}
	}

	return nil
}
