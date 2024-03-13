package types

import (
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/out"

	"github.com/jedib0t/go-pretty/v6/table"
)

//go:generate mockery --name=ISourceSettings --structname=ISourceSettingsMock --filename=series_mock.go --inpackage
type ISourceSettings interface {
	Print(cmd ICommand) error
}

type Series struct {
	Name           string          `json:"name"            mapstructure:"name"            yaml:"name"`
	ID             string          `json:"id"              mapstructure:"id"              yaml:"id"`
	ISBNBlacklist  []string        `json:"isbn_blacklist"  mapstructure:"isbn_blacklist"  yaml:"isbn_blacklist"`
	Source         SourceType      `json:"source"          mapstructure:"source"          yaml:"source"`
	SourceSettings ISourceSettings `json:"source_settings" mapstructure:"source_settings" yaml:"source_settings"`
	Key            string          `json:"key"             mapstructure:"key"             yaml:"key"`
}

func (s Series) Print(cmd ICommand) error {
	t := out.NewTable(cmd.OutOrStdout())
	t.AppendHeader(table.Row{"Key", "Name", "Source", "ID"})
	t.AppendRow(table.Row{s.Key, s.Name, s.Source, s.ID})
	t.Render()

	if len(s.ISBNBlacklist) > 0 {
		t := out.NewTable(cmd.OutOrStdout())
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

func (s Series) String() string {
	return fmt.Sprintf("%s (%s)", s.Name, s.Source)
}

func GetSetting[t ISourceSettings](s Series) (empty t, err error) {
	if s.SourceSettings == nil {
		return empty, err
	}

	settings, ok := s.SourceSettings.(t)
	if ok {
		return settings, nil
	}

	return empty, fmt.Errorf("setting type not correct")
}
