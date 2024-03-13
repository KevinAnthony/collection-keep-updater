package yen

import (
	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/jedib0t/go-pretty/v6/table"
)

type yenSettings struct{}

func (y yenSettings) Print(cmd types.ICommand) error {
	t := out.NewTable(cmd.OutOrStdout())
	t.AppendHeader(table.Row{"No Yen Settings"})

	t.Render()

	return nil
}
