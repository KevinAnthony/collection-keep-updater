package yen

import (
	"github.com/kevinanthony/collection-keep-updater/out"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type yenSettings struct{}

func (y yenSettings) Print(cmd *cobra.Command) error {
	t := out.NewTable(cmd)
	t.AppendHeader(table.Row{"No Yen Settings"})

	t.Render()

	return nil
}
