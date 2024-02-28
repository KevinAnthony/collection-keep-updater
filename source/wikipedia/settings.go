package wikipedia

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kevinanthony/collection-keep-updater/out"

	"github.com/spf13/cobra"
)

const (
	volumeF = "volume-header"
	titleF  = "title-header"
	isbnF   = "isbn-header"
	tableF  = "table-numbers"
)

var (
	volumeV string
	titleV  string
	isbnV   string
	tableV  []int
)

type wikiSettings struct {
	VolumeHeader *string `json:"volume_header" yaml:"volume_header"`
	TitleHeader  *string `json:"title_header"  yaml:"title_header"`
	ISBNHeader   *string `json:"isbn_header"   yaml:"isbn_header"`
	Table        []int   `json:"tables"        yaml:"tables"`
}

func (w wikiSettings) Print(cmd *cobra.Command) error {
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
