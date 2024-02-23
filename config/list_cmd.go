package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list one or all configurations",
	RunE:  RunList,
	Args:  ValidateListArgs,
}

func ValidateListArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		if series == library {
			return errors.New("configuration name given but no single configuration type set")
		}
	}

	return nil
}

func RunList(cmd *cobra.Command, args []string) error {
	var settingsKey string
	if len(args) > 0 {
		settingsKey = args[0]
	}

	cfg, err := ctxu.GetConfigCtx(cmd)
	if err != nil {
		return err
	}

	switch {
	case series:
		err = printSeries(cmd, cfg, settingsKey)
	case library:
		err = printLibrary(cmd, cfg, settingsKey)
	default:
		err = printSeriesBasic(cmd, cfg)
		if err != nil {
			return err
		}
		err = printLibraryBasic(cmd, cfg)
	}
	if err != nil {
		return err
	}

	return nil
}

func printSeries(cmd *cobra.Command, cfg types.Config, key string) error {
	if len(key) == 0 {
		return printSeriesBasic(cmd, cfg)
	}
	for _, s := range cfg.Series {
		if s.Key == key {
			return s.Print(cmd)
		}
	}

	return fmt.Errorf("key: %s not found in series configuration", key)
}

func printSeriesBasic(cmd *cobra.Command, cfg types.Config) error {
	t := out.NewTable(cmd)
	t.AppendHeader(table.Row{"Series"})
	t.AppendHeader(table.Row{"Key", "Name", "ID", "Source Type"})

	for _, s := range cfg.Series {
		t.AppendRow(table.Row{s.Key, s.Name, s.ID, s.Source})
	}

	t.Render()

	return nil
}

func printLibrary(cmd *cobra.Command, cfg types.Config, key string) error {
	if len(key) == 0 {
		return printLibraryBasic(cmd, cfg)
	}
	t := out.NewTable(cmd)
	t.AppendHeader(table.Row{"Library", "Wanted Collection ID", "Other Collection IDs", "API Key"}, out.AutoMergeRow)

	for _, l := range cfg.Libraries {
		if string(l.Name) == key {
			t.AppendRow(table.Row{l.Name, l.WantedColID, strings.Join(l.OtherColIDs, ","), out.Partial(l.APIKey, 50)})

			t.Render()

			return nil
		}
	}

	return fmt.Errorf("name: %s not found in library configuration", key)
}

func printLibraryBasic(cmd *cobra.Command, cfg types.Config) error {
	t := out.NewTable(cmd)
	t.AppendHeader(table.Row{"Library", "Wanted Collection ID", "Other Collection IDs", "API Key"}, out.AutoMergeRow)

	for _, l := range cfg.Libraries {
		t.AppendRow(table.Row{l.Name, l.WantedColID, strings.Join(l.OtherColIDs, ","), out.Partial(l.APIKey, 50)})
	}

	t.Render()

	return nil
}
