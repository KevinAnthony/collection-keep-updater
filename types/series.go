package types

import (
	"reflect"

	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ISourceSettings interface {
	Print(cmd *cobra.Command) error
}

type Series struct {
	Name           string          `json:"name"            mapstructure:"name"`
	ID             string          `json:"id"              mapstructure:"id"`
	ISBNBlacklist  []string        `json:"isbn_blacklist"  mapstructure:"isbn_blacklist"`
	Source         SourceType      `json:"source"          mapstructure:"source"`
	SourceSettings ISourceSettings `json:"source_settings" mapstructure:"source_settings"`
	Key            string          `json:"key"             mapstructure:"key"`
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

const (
	seriesNameFlag   = "name"
	seriesKeyFlag    = "key"
	seriesIDFlag     = "url"
	seriesSourceFlag = "source"
	seriesBlacklist  = "blacklist"
)

func SeriesSetFlags(cmd *cobra.Command) {
	//WikipediaSetFlags(cmd)
	//VisSetFlags(cmd)

	cmd.PersistentFlags().String(seriesNameFlag, "", "name of the series.")
	cmd.PersistentFlags().String(seriesKeyFlag, "", "unique key of the series.")
	cmd.PersistentFlags().String(seriesIDFlag, "", "url to be parsed for the series, extracting the ID.")
	cmd.PersistentFlags().String(seriesSourceFlag, "", "type of source to be added. [viz,wikipieda]")
	cmd.PersistentFlags().StringArray(seriesBlacklist, []string{}, "name of the series.")
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
