package yen

import (
	"errors"
	"strings"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/spf13/cobra"
)

type settingsHelper struct {
}

func (s settingsHelper) SourceSettingFromConfig(_ map[string]interface{}) types.ISourceSettings {
	return &yenSettings{}
}

func (s settingsHelper) SourceSettingFromFlags(_ *cobra.Command, sourceSetting types.ISourceSettings) (types.ISourceSettings, error) {
	settings, ok := sourceSetting.(*yenSettings)
	if !ok {
		settings = &yenSettings{}
	}

	return settings, nil
}

func (s settingsHelper) SetFlags(_ *cobra.Command) {
}

func (s settingsHelper) GetIDFromURL(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("unknown/unset url.  url is required")
	}
	if !strings.HasPrefix(url, "https://yenpress.com/series/") {
		return "", errors.New("url is malformed")
	}

	return strings.TrimPrefix(url, "https://yenpress.com/series/"), nil
}
