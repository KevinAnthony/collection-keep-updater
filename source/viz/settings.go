package viz

import (
	"time"

	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/jedib0t/go-pretty/v6/table"
)

type vizSettings struct {
	MaximumBacklog *int           `json:"maximum_backlog" yaml:"maximum_backlog"`
	Delay          *time.Duration `json:"delay_between"   yaml:"delay_between"`
}

func (v vizSettings) Print(cmd types.ICommand) error {
	t := out.NewTable(cmd.OutOrStdout())
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
