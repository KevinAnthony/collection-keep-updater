package out

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var AutoMergeRow = table.RowConfig{AutoMerge: true}

func NewTable(cmd *cobra.Command) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(cmd.OutOrStdout())
	t.SetStyle(table.StyleLight)

	return t
}
