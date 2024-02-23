package types

import (
	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type WikipediaSettings struct {
	Volume          *string `json:"volume_column_title"`
	Title           *string `json:"title_column_title"`
	ISBNColumnTitle *string `json:"isbn_column_title"`
	Table           []int   `json:"tables"`
}

func (w WikipediaSettings) Print(cmd *cobra.Command) error {
	t := out.NewTable(cmd)
	t.AppendHeader(table.Row{"Volume Column Name", "Title Column Name", "ISBN Column Title", "Tables On Page"})
	t.AppendRow(
		table.Row{
			out.ValueOrEmpty(w.Volume),
			out.ValueOrEmpty(w.Title),
			out.ValueOrEmpty(w.ISBNColumnTitle),
			out.IntSliceToStrOrEmpty(w.Table),
		},
	)

	t.Render()

	return nil
}

func newWikipediaSettings(data map[string]interface{}) *WikipediaSettings {
	if len(data) == 0 {
		return nil
	}

	settings := WikipediaSettings{
		Volume:          utils.GetPtr[string](data, "volume_column_title"),
		Title:           utils.GetPtr[string](data, "title_column_title"),
		ISBNColumnTitle: utils.GetPtr[string](data, "isbn_column_title"),
		Table:           utils.GetArray[int](data, "tables"),
	}

	if settings.ISBNColumnTitle == nil || len(*settings.ISBNColumnTitle) == 0 {
		return nil
	}

	return &settings
}
