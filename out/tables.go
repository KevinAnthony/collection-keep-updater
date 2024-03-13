package out

import (
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
)

var AutoMergeRow = table.RowConfig{AutoMerge: true}

func NewTable(writer io.Writer) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(writer)
	t.SetStyle(table.StyleLight)

	return t
}
