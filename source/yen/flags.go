package yen

import (
	"errors"
	"strings"

	"github.com/kevinanthony/collection-keep-updater/source"
	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/spf13/cobra"
)

func init() {
	source.RegisterConfigCallbacks(types.YenSource, &source.ConfigCallback{
		SetFlagsFunc:                setFlags,
		SourceSettingFromConfigFunc: newYenSettings,
		SourceSettingFromFlagsFunc:  configFromFlags,
		GetIDFromURL:                parseURLToID,
	})
}

func newYenSettings(_ map[string]interface{}) types.ISourceSettings {
	return &yenSettings{}
}

func configFromFlags(_ *cobra.Command, sourceSetting types.ISourceSettings) (types.ISourceSettings, error) {
	settings, ok := sourceSetting.(*yenSettings)
	if !ok {
		settings = &yenSettings{}
	}

	return settings, nil
}

func setFlags(_ *cobra.Command) {

}

func parseURLToID(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("unknown/unset url.  url is required")
	}
	if !strings.HasPrefix(url, "https://yenpress.com/series/") {
		return "", errors.New("url is malformed")
	}

	return strings.TrimPrefix(url, "https://yenpress.com/series/"), nil
}
