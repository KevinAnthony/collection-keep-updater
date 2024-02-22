package types

import (
	"time"

	"github.com/kevinanthony/collection-keep-updater/utils"
)

type VizSettings struct {
	MaximumBacklog *int           `json:"maximum_backlog"`
	Delay          *time.Duration `json:"delay_between"`
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
