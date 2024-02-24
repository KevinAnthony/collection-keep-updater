package types

import (
	"errors"
	"reflect"

	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	seriesNameFlag      = "name"
	seriesKeyFlag       = "key"
	seriesURLFlag       = "url"
	seriesSourceFlag    = "source"
	seriesBlacklistFlag = "blacklist"
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

	cmd.Flags().StringVar(&seriesNameV, seriesNameFlag, "", "name of the series.")
	cmd.Flags().StringVar(&seriesKeyV, seriesKeyFlag, "", "unique key of the series.")
	cmd.Flags().StringVar(&seriesURLV, seriesURLFlag, "", "url to be parsed for the series, extracting the ID.")
	cmd.Flags().StringVar(&seriesSourceV, seriesSourceFlag, "", "type of source to be added. [viz, wikipieda]")
	cmd.Flags().StringArrayVar(&seriesBlacklistV, seriesBlacklistFlag, []string{}, "list of ISBNs to be ignored.")
}

func NewSeriesConfig(cmd *cobra.Command) (Series, error) {
	return SeriesConfigFromFlags(cmd, Series{})
}

func SeriesConfigFromFlags(cmd *cobra.Command, series Series) (Series, error) {
	series.Name = getFlagOrNil[string](cmd, seriesNameFlag, seriesNameV)
	series.Key = getFlagOrNil[string](cmd, seriesKeyFlag, seriesKeyV)
	series.Source = SourceType(getFlagOrNil[string](cmd, seriesSourceFlag, seriesSourceV))
	series.ISBNBlacklist = getFlagOrNil[[]string](cmd, seriesBlacklistFlag, seriesBlacklistV)

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
