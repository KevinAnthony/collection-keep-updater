package types

import (
	"github.com/kevinanthony/collection-keep-updater/out"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
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
