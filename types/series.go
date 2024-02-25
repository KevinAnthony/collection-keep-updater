package types

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	seriesNameF      = "name"
	seriesKeyF       = "key"
	seriesURLF       = "url"
	seriesSourceF    = "source"
	seriesBlacklistF = "blacklist"
)

type ISourceSettings interface {
	Print(cmd *cobra.Command) error
}

type Series struct {
	Name           string          `json:"name"            mapstructure:"name"            yaml:"name"`
	ID             string          `json:"id"              mapstructure:"id"              yaml:"id"`
	ISBNBlacklist  []string        `json:"isbn_blacklist"  mapstructure:"isbn_blacklist"  yaml:"isbn_blacklist"`
	Source         SourceType      `json:"source"          mapstructure:"source"          yaml:"source"`
	SourceSettings ISourceSettings `json:"source_settings" mapstructure:"source_settings" yaml:"source_settings"`
	Key            string          `json:"key"             mapstructure:"key"             yaml:"key"`
}

func (s Series) Print(cmd *cobra.Command) error {
	t := out.NewTable(cmd)
	t.AppendHeader(table.Row{"Key", "Name", "Source", "ID"})
	t.AppendRow(table.Row{s.Key, s.Name, s.Source, s.ID})
	t.Render()

	if len(s.ISBNBlacklist) > 0 {
		t := out.NewTable(cmd)
		t.AppendHeader(table.Row{"ISBN Blacklist"})

		for _, isbn := range s.ISBNBlacklist {
			t.AppendRow(table.Row{isbn})
		}

		t.Render()
	}

	if s.SourceSettings != nil {
		s.SourceSettings.Print(cmd)
	}

	return nil
}

var (
	seriesNameV      string
	seriesKeyV       string
	seriesURLV       string
	seriesSourceV    string
	seriesBlacklistV []string
)

func SeriesSetFlags(cmd *cobra.Command) {
	wikipediaSetFlags(cmd)
	vizSetFlags(cmd)

	cmd.Flags().StringVar(&seriesNameV, seriesNameF, "", "name of the series.")
	cmd.Flags().StringVar(&seriesKeyV, seriesKeyF, "", "unique key of the series.")
	cmd.Flags().StringVar(&seriesURLV, seriesURLF, "", "url to be parsed for the series, extracting the ID.")
	cmd.Flags().StringVar(&seriesSourceV, seriesSourceF, "", "type of source to be added. [viz, wikipieda]")
	cmd.Flags().StringArrayVar(&seriesBlacklistV, seriesBlacklistF, []string{}, "list of ISBNs to be ignored.")
}

func NewSeriesConfig(cmd *cobra.Command) (Series, error) {
	return seriesConfigFromFlags(cmd, Series{})
}

func EditSeries(cmd *cobra.Command, cfg Config) (*Series, error) {
	key := getFlagOrDefault[string](cmd, seriesKeyF, seriesKeyV, "")
	if len(key) == 0 {
		return nil, errors.New("key flag is required for edit")
	}

	var series *Series
	for _, s := range cfg.Series {
		if s.Key == key {
			series = &s

			break
		}
	}
	//
	if series == nil {
		return nil, fmt.Errorf("edit: series key %s not found in config", key)
	}

	s, err := seriesConfigFromFlags(cmd, *series)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func seriesConfigFromFlags(cmd *cobra.Command, series Series) (Series, error) {
	series.Name = getFlagOrDefault[string](cmd, seriesNameF, seriesNameV, series.Name)
	series.Key = getFlagOrDefault[string](cmd, seriesKeyF, seriesKeyV, series.Key)
	series.Source = SourceType(getFlagOrDefault[string](cmd, seriesSourceF, seriesSourceV, string(series.Source)))
	series.ISBNBlacklist = getFlagOrDefault[[]string](cmd, seriesBlacklistF, seriesBlacklistV, series.ISBNBlacklist)

	switch series.Source {
	case WikipediaSource:
		return wikiConfigFromFlags(cmd, series)
	case VizSource:
		return vizConfigFromFlags(cmd, series)
	default:
		return series, errors.New("unknown/unset source.  source is required")
	}
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
			Key:           utils.Get[string](values, "key"),
			Source:        SourceType(utils.Get[string](values, "source")),
			ISBNBlacklist: utils.GetArray[string](values, "isbn_blacklist"),
		}

		series.SourceSettings = getSetting(series.Source, sourceValue)

		return series, nil
	},
	)
}

func getSetting(source SourceType, data map[string]interface{}) ISourceSettings {
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
