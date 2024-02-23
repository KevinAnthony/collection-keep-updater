package types

import (
	"time"

	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type VizSettings struct {
	MaximumBacklog *int           `json:"maximum_backlog"`
	Delay          *time.Duration `json:"delay_between"`
}

func (v VizSettings) Print(cmd *cobra.Command) error {
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

func newVizSettings(data map[string]interface{}) *VizSettings {
	if len(data) == 0 {
		return nil
	}

	settings := VizSettings{
		MaximumBacklog: utils.GetPtr[int](data, "maximum_backlog"),
	}

	delayStr := utils.Get[string](data, "delay_between")
	if len(delayStr) > 0 {
		delay, err := time.ParseDuration(delayStr)
		if err != nil {
			// TODO: log error
		} else {
			settings.Delay = &delay
		}
	}

	return &settings
}
