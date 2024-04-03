package viz

import (
	"strings"
	"time"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/pkg/errors"
)

const (
	maxBacklogF = "viz-max-backlog"
	getDelayF   = "viz-get-delay"
)

type settingsHelper struct{}

func SetFlags(cmd types.ICommand) {
	cmd.PersistentFlags().Int(maxBacklogF, 0, "how many volumes from the end to check.")
	cmd.PersistentFlags().String(getDelayF, "", "how long a delay to wait between each request, in go time.Duration format.")
}

func (v settingsHelper) SourceSettingFromFlags(cmd types.ICommand, sourceSetting types.ISourceSettings) (types.ISourceSettings, error) {
	settings, ok := sourceSetting.(*vizSettings)
	if !ok {
		settings = &vizSettings{}
	}

	settings.MaximumBacklog = utils.GetFlagIntPtr(cmd, maxBacklogF)

	delayStr := utils.GetFlagStringPtr(cmd, getDelayF)
	if delayStr != nil && len(*delayStr) > 0 {
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
