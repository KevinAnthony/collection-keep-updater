package config

import (
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a configuration",
	RunE:  runRemove,
	Args:  validateRemoveArgs,
}

func validateRemoveArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("remove: key not provided")
	}

	return nil
}

func runRemove(cmd *cobra.Command, args []string) error {
	cfg, err := ctxu.GetConfig(cmd)
	if err != nil {
		return err
	}

	settingsKey := args[0]

	switch {
	case isSeries:
		return removeSeries(cfg, settingsKey)
	case isLibrary:
		return removeLibrary(cfg, settingsKey)
	default:
		return errors.New("unknown configuration type")
	}
}

func removeSeries(cfg types.Config, key string) error {
	for i, s := range cfg.Series {
		if s.Key == key {
			cfg.Series = append(cfg.Series[:i], cfg.Series[i+1:]...)

			viper.Set("series", cfg.Series)

			return viper.WriteConfig()
		}
	}

	return fmt.Errorf("remove: unknown series: %s", key)
}

func removeLibrary(cfg types.Config, key string) error {
	for i, l := range cfg.Libraries {
		if string(l.Name) == key {
			cfg.Libraries = append(cfg.Libraries[:i], cfg.Libraries[i+1:]...)

			viper.Set("libraries", cfg.Libraries)

			return viper.WriteConfig()
		}
	}

	return fmt.Errorf("remove: unknown library: %s", key)
}
