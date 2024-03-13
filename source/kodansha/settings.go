package kodansha

import (
	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/jedib0t/go-pretty/v6/table"
)

type kondashaSettings struct{}

func (k kondashaSettings) Print(cmd types.ICommand) error {
	t := out.NewTable(cmd.OutOrStdout())
	t.AppendHeader(table.Row{"No Kondasha Settings"})

	t.Render()

	return nil
}
