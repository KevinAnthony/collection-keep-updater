package config

import (
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit an existing configuration",
	RunE:  runEdit,
}

func init() {
	types.SeriesSetFlags(editCmd)
	editCmd.PersistentFlags().BoolVarP(&try, "test-config", "t", false, "test the configuration by calling source and outputting result.")
	editCmd.PersistentFlags().BoolVarP(&write, "write-config", "w", false, "save the configuration.")

	editCmd.MarkFlagsOneRequired("test-config", "write-config")
	editCmd.MarkFlagsMutuallyExclusive("test-config", "write-config")
}

func runEdit(cmd *cobra.Command, args []string) error {
	cfg, err := ctxu.GetConfig(cmd)
	if err != nil {
		return err
	}

	switch {
	case isSeries:
		s, err := types.EditSeries(cmd, cfg)
		if err != nil {
			return err
		}

		if try {
			sources, err := ctxu.GetSources(cmd)
			if err != nil {
				return err
			}

			source, found := sources[s.Source]
			if !found {
				return fmt.Errorf("source type %s not found in source map", s.Source)
			}

			books, err := source.GetISBNs(cmd.Context(), *s)
			if err != nil {
				return err
			}

			t := out.NewTable(cmd)
			t.AppendHeader(table.Row{"Title", "Volume", "ISBN 10", "ISBN 13"})
			for _, book := range books {
				t.AppendRow(table.Row{book.Title, book.Volume, book.ISBN10, book.ISBN13})
			}

			t.Render()
		}

		if write {
			cfg, err := ctxu.GetConfig(cmd)
			if err != nil {
				return err
			}

			cfg.Series = append(cfg.Series, *s)
			viper.Set("series", cfg.Series)

			return viper.WriteConfig()
		}
	}

	return nil
}
