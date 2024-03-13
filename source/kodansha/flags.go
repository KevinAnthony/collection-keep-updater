package kodansha

import (
	"strings"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/pkg/errors"
)

type settingsHelper struct{}

func (s settingsHelper) SourceSettingFromConfig(_ map[string]interface{}) types.ISourceSettings {
	return kondashaSettings{}
}

func (s settingsHelper) SourceSettingFromFlags(_ types.ICommand, sourceSetting types.ISourceSettings) (types.ISourceSettings, error) {
	settings, ok := sourceSetting.(*kondashaSettings)
	if !ok {
		settings = &kondashaSettings{}
	}

	return settings, nil
}

func (s settingsHelper) GetIDFromURL(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("unknown/unset url.  url is required")
	}

	if !strings.HasPrefix(url, baseURL+seriesSlug) {
		return "", errors.New("url is malformed")
	}

	return strings.TrimPrefix(url, baseURL+seriesSlug), nil
}

func SetFlags(_ types.ICommand) {
}
