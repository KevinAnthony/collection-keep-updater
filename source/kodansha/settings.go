package kodansha

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/spf13/cobra"
)

type kondashaSettings struct{}

func (k kondashaSettings) Print(cmd *cobra.Command) error {
	t := out.NewTable(cmd)
	t.AppendHeader(table.Row{"No Kondasha Settings"})

	t.Render()

	return nil
}
