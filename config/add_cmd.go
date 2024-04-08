package config

import (
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/spf13/cobra"
)

const (
	testF  = "test-config"
	writeF = "write-config"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new configuration",
	RunE:  types.CmdRunE(runAdd),
}

func init() {
	addCmd.PersistentFlags().BoolP(testF, "t", false, "test the configuration by calling source and outputting result.")
	addCmd.PersistentFlags().BoolP(writeF, "w", false, "save the configuration.")

	addCmd.MarkFlagsOneRequired("test-config", "write-config")
	addCmd.MarkFlagsMutuallyExclusive("test-config", "write-config")

	seriesSetFlags(addCmd)
}

func runAdd(cmd types.ICommand, args []string) error {
	switch {
	case utils.GetFlagBool(cmd, seriesFlag):
		s, err := newSeriesConfig(cmd)
		if err != nil {
			return err
		}

		if utils.GetFlagBool(cmd, testF) {
			source, err := ctxu.GetSource(cmd, s.Source)
			if err != nil {
				return err
			}

			books, err := source.GetISBNs(cmd.Context(), s)
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

			cfg.Series = append(cfg.Series, s)
			icfg := ctxu.GetConfigReader(cmd)

			icfg.Set("series", cfg.Series)

			return icfg.WriteConfig()
		}
	}

	return nil
}
