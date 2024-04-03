package wikipedia

import (
	"strings"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/pkg/errors"
)

const (
	volumeF = "wiki-volume-header"
	titleF  = "wiki-title-header"
	isbnF   = "wiki-isbn-header"
	tableF  = "wiki-table-numbers"
)

type settingsHelper struct{}

func (s settingsHelper) GetIDFromURL(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("unknown/unset url. url is required")
	}
	if !strings.HasPrefix(url, "https://en.wikipedia.org/wiki/") {
		return "", errors.New("url is malformed")
	}

	return strings.TrimPrefix(url, "https://en.wikipedia.org/wiki/"), nil
}

func (s settingsHelper) SourceSettingFromFlags(cmd types.ICommand, sourceSetting types.ISourceSettings) (types.ISourceSettings, error) {
	settings, ok := sourceSetting.(*wikiSettings)
	if !ok {
		settings = &wikiSettings{}
	}

	settings.VolumeHeader = utils.GetFlagStringPtr(cmd, volumeF)
	settings.TitleHeader = utils.GetFlagStringPtr(cmd, titleF)
	settings.ISBNHeader = utils.GetFlagStringPtr(cmd, isbnF)
	settings.Table = utils.GetFlagIntSlice(cmd, tableF)

	return settings, nil
}

func (s settingsHelper) SourceSettingFromConfig(data map[string]interface{}) types.ISourceSettings {
	if len(data) == 0 {
		return wikiSettings{}
	}

	settings := wikiSettings{
		VolumeHeader: utils.GetPtr[string](data, "volume_header"),
		TitleHeader:  utils.GetPtr[string](data, "title_header"),
		ISBNHeader:   utils.GetPtr[string](data, "isbn_header"),
		Table:        utils.GetArray[int](data, "tables"),
	}

	if settings.ISBNHeader == nil || len(*settings.ISBNHeader) == 0 {
		return wikiSettings{}
	}

	return settings
}

func SetFlags(cmd types.ICommand) {
	cmd.PersistentFlags().String(volumeF, "", "header of the column that has the volume number.")
	cmd.PersistentFlags().String(titleF, "", "header of the column that has the title.")
	cmd.PersistentFlags().String(isbnF, "", "header of the column that has the ISBN number(required).")
	cmd.PersistentFlags().IntSlice(tableF, []int{}, "tables to include, zero indexed. skip for all tables.")
}
