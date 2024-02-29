package viz

import (
	"strings"
	"time"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	maxBacklogF = "viz-max-backlog"
	getDelayF   = "viz-get-delay"
)

type vizFlags struct {
	maxBacklogV int
	getDelayV   string
}

func (v vizFlags) SetFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVar(&v.maxBacklogV, maxBacklogF, 0, "how many volumes from the end to check.")
	cmd.PersistentFlags().StringVar(&v.getDelayV, getDelayF, "", "how long a delay to wait between each request, in go time.Duration format.")
}

func (v vizFlags) SourceSettingFromFlags(cmd *cobra.Command, sourceSetting types.ISourceSettings) (types.ISourceSettings, error) {
	settings, ok := sourceSetting.(*vizSettings)
	if !ok {
		settings = &vizSettings{}
	}

	settings.MaximumBacklog = utils.GetFlagOrDefault[*int](cmd, maxBacklogF, &v.maxBacklogV, settings.MaximumBacklog)

	var str string
	if settings.Delay != nil {
		str = settings.Delay.String()
	}

	delayStr := utils.GetFlagOrDefault[*string](cmd, getDelayF, &v.getDelayV, &str)
	if delayStr != nil {
		delay, err := time.ParseDuration(*delayStr)
		if err != nil {
			return settings, errors.Wrap(err, "viz: cannot parse delay "+*delayStr)
		}
		settings.Delay = &delay
	}

	return settings, nil
}

func (v vizFlags) GetIDFromURL(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("unknown/unset url. url is required")
	}

	// regex maybe?
	if !strings.HasPrefix(url, "https://www.viz.com/read/manga/") || !strings.HasSuffix(url, "/all") {
		return "", errors.New("url is malformed")
	}

	url = strings.TrimPrefix(url, "https://www.viz.com/read/manga/")
	return strings.TrimSuffix(url, "/all"), nil
}

func (v vizFlags) SourceSettingFromConfig(data map[string]interface{}) types.ISourceSettings {
	if len(data) == 0 {
		return nil
	}

	settings := vizSettings{
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
