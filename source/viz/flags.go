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

var (
	maxBacklogV int
	getDelayV   string
)

type settingsHelper struct{}

func SetFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVar(&maxBacklogV, maxBacklogF, 0, "how many volumes from the end to check.")
	cmd.PersistentFlags().StringVar(&getDelayV, getDelayF, "", "how long a delay to wait between each request, in go time.Duration format.")
}

func (v settingsHelper) SourceSettingFromFlags(cmd *cobra.Command, sourceSetting types.ISourceSettings) (types.ISourceSettings, error) {
	settings, ok := sourceSetting.(*vizSettings)
	if !ok {
		settings = &vizSettings{}
	}

	settings.MaximumBacklog = utils.GetFlagOrDefault[*int](cmd, maxBacklogF, &maxBacklogV, settings.MaximumBacklog)

	var str string
	if settings.Delay != nil {
		str = settings.Delay.String()
	}

	delayStr := utils.GetFlagOrDefault[*string](cmd, getDelayF, &getDelayV, &str)
	if delayStr != nil {
		delay, err := time.ParseDuration(*delayStr)
		if err != nil {
			return settings, errors.Wrap(err, "viz: cannot parse delay "+*delayStr)
		}
		settings.Delay = &delay
	}

	return settings, nil
}

func (v settingsHelper) GetIDFromURL(url string) (string, error) {
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

func (v settingsHelper) SourceSettingFromConfig(data map[string]interface{}) types.ISourceSettings {
	if len(data) == 0 {
		return vizSettings{}
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

	return settings
}
