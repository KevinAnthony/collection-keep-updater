package viz

import (
	"time"

	"github.com/kevinanthony/collection-keep-updater/out"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

const (
	maxBacklogF = "max-backlog"
	getDelayF   = "get-delay"
)

var (
	maxBacklogV int
	getDelayV   string
)

type vizSettings struct {
	MaximumBacklog *int           `json:"maximum_backlog" yaml:"maximum_backlog"`
	Delay          *time.Duration `json:"delay_between"   yaml:"delay_between"`
}

func (v vizSettings) Print(cmd *cobra.Command) error {
	t := out.NewTable(cmd)
	t.AppendHeader(table.Row{"Maximum Backlog", "Delay"})
	t.AppendRow(
		table.Row{
			out.ValueOrEmpty(v.MaximumBacklog),
			out.ValueOrEmpty(v.Delay),
		},
	)
	t.Render()

	return nil
}
