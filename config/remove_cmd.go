package config

import (
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a configuration",
	RunE:  types.CmdRunArgs(runRemove),
	Args:  types.CmdRunArgs(validateRemoveArgs),
}

func validateRemoveArgs(_ types.ICommand, args []string) error {
	if len(args) == 0 {
		return errors.New("remove: key not provided")
	}

	return nil
}

func runRemove(cmd types.ICommand, args []string) error {
	cfg, err := ctxu.GetConfig(cmd)
	if err != nil {
		return err
	}

	settingsKey := args[0]

	switch {
	case utils.GetFlagBool(cmd, seriesFlag):
		return removeSeries(cmd, cfg, settingsKey)
	case utils.GetFlagBool(cmd, libraryFlag):
		return removeLibrary(cmd, cfg, settingsKey)
	default:
		return errors.New("unknown configuration type")
	}
}

func removeSeries(cmd types.ICommand, cfg types.Config, key string) error {
	for i, s := range cfg.Series {
		if s.Key == key {
			cfg.Series = append(cfg.Series[:i], cfg.Series[i+1:]...)

			icfg := ctxu.GetConfigReader(cmd)

			icfg.Set("series", cfg.Series)

			return icfg.WriteConfig()
		}
	}

	return fmt.Errorf("remove: unknown series: %s", key)
}

func removeLibrary(cmd types.ICommand, cfg types.Config, key string) error {
	for i, l := range cfg.Libraries {
		if string(l.Name) == key {
			cfg.Libraries = append(cfg.Libraries[:i], cfg.Libraries[i+1:]...)

			icfg := ctxu.GetConfigReader(cmd)

			icfg.Set("libraries", cfg.Libraries)

			return icfg.WriteConfig()
		}
	}

	return fmt.Errorf("remove: unknown library: %s", key)
}
