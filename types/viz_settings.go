package types

import (
	"strings"
	"time"

	"github.com/kevinanthony/collection-keep-updater/out"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	vizMaxBacklogF = "max-backlog"
	vizGetDelayF   = "get-delay"
)

var (
	vizMaxBacklogV int
	vizGetDelayV   string
)

type VizSettings struct {
	MaximumBacklog *int           `json:"maximum_backlog" yaml:"maximum_backlog"`
	Delay          *time.Duration `json:"delay_between"   yaml:"delay_between"`
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

func vizSetFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVar(&vizMaxBacklogV, vizMaxBacklogF, 0, "how many volumes from the end to check.")
	cmd.PersistentFlags().StringVar(&vizGetDelayV, vizGetDelayF, "", "how long a delay to wait between each request, in go time.Duration format.")
}

func vizConfigFromFlags(cmd *cobra.Command, series Series) (Series, error) {
	settings, ok := series.SourceSettings.(VizSettings)
	if !ok {
		settings = VizSettings{}
	}

	url := getFlagOrNil[string](cmd, seriesURLFlag, seriesURLV)
	if len(url) == 0 {
		return series, errors.New("unknown/unset url. url is required")
	}

	// regex maybe?
	if !strings.HasPrefix(url, "https://www.viz.com/read/manga/") || !strings.HasSuffix(url, "/all") {
		return series, errors.New("url is malformed")
	}
	url = strings.TrimPrefix(url, "https://www.viz.com/read/manga/")
	series.ID = strings.TrimSuffix(url, "/all")

	settings.MaximumBacklog = getFlagOrNil[*int](cmd, vizMaxBacklogF, &vizMaxBacklogV)
	delayStr := getFlagOrNil[*string](cmd, vizGetDelayF, &vizGetDelayV)
	if delayStr != nil {
		delay, err := time.ParseDuration(*delayStr)
		if err != nil {
			return series, errors.Wrap(err, "viz: cannot parse delay "+*delayStr)
		}
		settings.Delay = &delay
	}
	series.SourceSettings = &settings

	return series, nil
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
