package types

import (
	"strings"

	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	wikipediaVolumeFlag = "volume-header"
	wikipediaTitleFlag  = "title-header"
	wikipediaISBNFlag   = "isbn-header"
	wikipediaTableFlag  = "table-numbers"
)

var (
	wikipediaVolumeV string
	wikipediaTitleV  string
	wikipediaISBNV   string
	wikipediaTableV  []int
)

type WikipediaSettings struct {
	VolumeHeader *string `json:"volume_header" yaml:"volume_header"`
	TitleHeader  *string `json:"title_header"  yaml:"title_header"`
	ISBNHeader   *string `json:"isbn_header"   yaml:"isbn_header"`
	Table        []int   `json:"tables"        yaml:"tables"`
}

func (w WikipediaSettings) Print(cmd *cobra.Command) error {
	t := out.NewTable(cmd)
	t.AppendHeader(table.Row{"Volume Column Name", "Title Column Name", "ISBN Column Title", "Tables On Page"})
	t.AppendRow(
		table.Row{
			out.ValueOrEmpty(w.VolumeHeader),
			out.ValueOrEmpty(w.TitleHeader),
			out.ValueOrEmpty(w.ISBNHeader),
			out.IntSliceToStrOrEmpty(w.Table),
		},
	)

	t.Render()

	return nil
}

func wikiConfigFromFlags(cmd *cobra.Command, series Series) (Series, error) {
	settings, ok := series.SourceSettings.(WikipediaSettings)
	if !ok {
		settings = WikipediaSettings{}
	}

	url := getFlagOrNil[string](cmd, seriesURLFlag, seriesURLV)
	if len(url) == 0 {
		return series, errors.New("unknown/unset url.  url is required")
	}
	if !strings.HasPrefix(url, "https://en.wikipedia.org/wiki/") {
		return series, errors.New("url is malformed")
	}

	series.ID = strings.TrimPrefix(url, "https://en.wikipedia.org/wiki/")

	settings.VolumeHeader = getFlagOrNil[*string](cmd, wikipediaVolumeFlag, &wikipediaVolumeV)
	settings.TitleHeader = getFlagOrNil[*string](cmd, wikipediaTitleFlag, &wikipediaTitleV)
	settings.ISBNHeader = getFlagOrNil[*string](cmd, wikipediaISBNFlag, &wikipediaISBNV)
	settings.Table = getFlagOrNil[[]int](cmd, wikipediaTableFlag, wikipediaTableV)

	series.SourceSettings = &settings

	return series, nil
}

func wikipediaSetFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&wikipediaVolumeV, wikipediaVolumeFlag, "", "header of the column that has the volume number.")
	cmd.PersistentFlags().StringVar(&wikipediaTitleV, wikipediaTitleFlag, "", "header of the column that has the title.")
	cmd.PersistentFlags().StringVar(&wikipediaISBNV, wikipediaISBNFlag, "", "header of the column that has the ISBN number(required).")
	cmd.PersistentFlags().IntSliceVar(&wikipediaTableV, wikipediaTableFlag, []int{}, "tables to include, zero indexed. skip for all tables.")
}

func newWikipediaSettings(data map[string]interface{}) *WikipediaSettings {
	if len(data) == 0 {
		return nil
	}

	settings := WikipediaSettings{
		VolumeHeader: utils.GetPtr[string](data, "volume_header"),
		TitleHeader:  utils.GetPtr[string](data, "title_header"),
		ISBNHeader:   utils.GetPtr[string](data, "isbn_header"),
		Table:        utils.GetArray[int](data, "tables"),
	}

	if settings.ISBNHeader == nil || len(*settings.ISBNHeader) == 0 {
		return nil
	}

	return &settings
}
