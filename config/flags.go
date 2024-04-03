package config

import (
	"fmt"
	"reflect"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/source/yen"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	seriesNameF      = "name"
	seriesKeyF       = "key"
	seriesURLF       = "url"
	seriesSourceF    = "source"
	seriesBlacklistF = "blacklist"
)

func seriesSetFlags(cmd types.ICommand) {
	cmd.PersistentFlags().String(seriesNameF, "", "name of the series.")
	cmd.PersistentFlags().String(seriesKeyF, "", "unique key of the series.")
	cmd.PersistentFlags().String(seriesURLF, "", "url to be parsed for the series, extracting the ID.")
	cmd.PersistentFlags().String(seriesSourceF, "", "type of source to be added. [viz, wikipieda]")
	cmd.PersistentFlags().StringArray(seriesBlacklistF, []string{}, "list of ISBNs to be ignored.")

	viz.SetFlags(cmd)
	wikipedia.SetFlags(cmd)
	yen.SetFlags(cmd)
}

func newSeriesConfig(cmd types.ICommand) (types.Series, error) {
	return seriesConfigFromFlags(cmd, types.Series{})
}

func editSeries(cmd types.ICommand, cfg types.Config) (*types.Series, error) {
	key := utils.GetFlagString(cmd, seriesKeyF)
	if len(key) == 0 {
		return nil, errors.New("key flag is required for edit")
	}

	var series *types.Series
	for _, s := range cfg.Series {
		if s.Key == key {
			series = &s

			break
		}
	}

	if series == nil {
		return nil, fmt.Errorf("edit: series key %s not found in config", key)
	}

	s, err := seriesConfigFromFlags(cmd, *series)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func seriesConfigFromFlags(cmd types.ICommand, series types.Series) (types.Series, error) {
	series.Name = utils.GetFlagString(cmd, seriesNameF)
	series.Key = utils.GetFlagString(cmd, seriesKeyF)
	series.Source = types.SourceType(utils.GetFlagString(cmd, seriesSourceF))
	series.ISBNBlacklist = utils.GetFlagStringSlice(cmd, seriesBlacklistF)

	url := utils.GetFlagString(cmd, seriesURLF)

	sourceSetting, err := ctxu.GetSourceSetting(cmd, series.Source)
	if err != nil {
		return series, errors.New("unknown/unset source.  source is required")
	}

	// if url is empty and series ID is already set, do not try and reset it
	if len(series.ID) == 0 || len(url) > 0 {
		if len(url) == 0 {
			return series, errors.New("unknown/unset url. url is required")
		}

		id, err := sourceSetting.GetIDFromURL(url)
		if err != nil {
			return series, err
		}

		series.ID = id
	}

	settings, err := sourceSetting.SourceSettingFromFlags(cmd, series.SourceSettings)
	if err != nil {
		return series, err
	}

	series.SourceSettings = settings

	return series, nil
}

func SeriesConfigHookFunc(cmd types.ICommand) viper.DecoderConfigOption {
	return viper.DecodeHook(func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}

		if t != reflect.TypeOf(types.Series{}) {
			return data, nil
		}

		values, ok := data.(map[string]interface{})
		if !ok {
			return data, nil
		}

		sourceValue := utils.Get[map[string]interface{}](values, "source_settings")
		series := types.Series{
			Name:          utils.Get[string](values, "name"),
			ID:            utils.Get[string](values, "id"),
			Key:           utils.Get[string](values, "key"),
			Source:        types.SourceType(utils.Get[string](values, "source")),
			ISBNBlacklist: utils.GetArray[string](values, "isbn_blacklist"),
		}

		series.SourceSettings = getSetting(cmd, series.Source, sourceValue)

		return series, nil
	},
	)
}

func getSetting(cmd types.ICommand, key types.SourceType, data map[string]interface{}) types.ISourceSettings {
	if len(data) == 0 {
		return nil
	}

	sourceSetting, err := ctxu.GetSourceSetting(cmd, key)
	if err != nil {
		return nil
	}

	return sourceSetting.SourceSettingFromConfig(data)
}
