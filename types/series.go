package types

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Series struct {
	Name           string     `json:"name"            mapstructure:"name"`
	ID             string     `json:"id"              mapstructure:"id"`
	ISBNBlacklist  []string   `json:"isbn_blacklist"  mapstructure:"isbn_blacklist"`
	Source         SourceType `json:"source"          mapstructure:"source"`
	SourceSettings any        `json:"source_settings" mapstructure:"source_settings"`
}

func (s *Series) UnmarshalJSON(data []byte) error {
	var raw map[string]*json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	name, err := utils.Unmarshal[string](raw["name"])
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal series name")
	}

	id, err := utils.Unmarshal[string](raw["id"])
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal series id")
	}

	source, err := utils.Unmarshal[SourceType](raw["source"])
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal series source")
	}

	blacklist, err := utils.Unmarshal[[]string](raw["isbn_blacklist"])
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal series blacklist")
	}

	s.Name = name
	s.ID = id
	s.Source = source
	s.ISBNBlacklist = blacklist

	switch s.Source {
	case WikipediaSource:
		settings, err := utils.Unmarshal[WikipediaSettings](raw["source_settings"])
		if err != nil {
			return errors.Wrap(err, "unable to unmarshal source_settings to wikipedia")
		}

		s.SourceSettings = settings
	case VizSource:
		settings, err := utils.Unmarshal[VizSettings](raw["source_settings"])
		if err != nil {
			return errors.Wrap(err, "unable to unmarshal source_settings to viz")
		}

		s.SourceSettings = settings
	default:
		return fmt.Errorf("unknown source type: %s", s.Source)
	}

	return nil
}

func SeriesConfigHookFunc() viper.DecoderConfigOption {
	return viper.DecodeHook(func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}

		if t != reflect.TypeOf(Series{}) {
			return data, nil
		}

		values, ok := data.(map[string]interface{})
		if !ok {
			return data, nil
		}

		sourceValue := utils.Get[map[string]interface{}](values, "source_settings")
		series := Series{
			Name:          utils.Get[string](values, "name"),
			ID:            utils.Get[string](values, "id"),
			Source:        SourceType(utils.Get[string](values, "source")),
			ISBNBlacklist: utils.GetArray[string](values, "source"),
		}

		series.SourceSettings = getSetting(series.Source, sourceValue)

		return series, nil
	},
	)
}

func getSetting(source SourceType, data map[string]interface{}) any {
	if len(data) == 0 {
		return nil
	}
	switch source {
	case WikipediaSource:
		return newWikipediaSettings(data)
	case VizSource:
		return newVizSettings(data)
	default:
		return nil
	}
}
