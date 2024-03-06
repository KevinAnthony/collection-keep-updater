package wikipedia

import (
	"strings"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	volumeF = "wiki-volume-header"
	titleF  = "wiki-title-header"
	isbnF   = "wiki-isbn-header"
	tableF  = "wiki-table-numbers"
)

var (
	volumeV string
	titleV  string
	isbnV   string
	tableV  []int
)

type settingsHelper struct{}

func (s settingsHelper) GetIDFromURL(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("unknown/unset url.  url is required")
	}
	if !strings.HasPrefix(url, "https://en.wikipedia.org/wiki/") {
		return "", errors.New("url is malformed")
	}

	return strings.TrimPrefix(url, "https://en.wikipedia.org/wiki/"), nil
}

func (s settingsHelper) SourceSettingFromFlags(cmd *cobra.Command, sourceSetting types.ISourceSettings) (types.ISourceSettings, error) {
	settings, ok := sourceSetting.(*wikiSettings)
	if !ok {
		settings = &wikiSettings{}
	}

	settings.VolumeHeader = utils.GetFlagOrDefault[*string](cmd, volumeF, &volumeV, settings.VolumeHeader)
	settings.TitleHeader = utils.GetFlagOrDefault[*string](cmd, titleF, &titleV, settings.TitleHeader)
	settings.ISBNHeader = utils.GetFlagOrDefault[*string](cmd, isbnF, &isbnV, settings.ISBNHeader)
	settings.Table = utils.GetFlagOrDefault[[]int](cmd, tableF, tableV, settings.Table)

	return settings, nil
}

func (s settingsHelper) SourceSettingFromConfig(data map[string]interface{}) types.ISourceSettings {
	if len(data) == 0 {
		return nil
	}

	settings := wikiSettings{
		VolumeHeader: utils.GetPtr[string](data, "volume_header"),
		TitleHeader:  utils.GetPtr[string](data, "title_header"),
		ISBNHeader:   utils.GetPtr[string](data, "isbn_header"),
		Table:        utils.GetArray[int](data, "tables"),
	}

	if settings.ISBNHeader == nil || len(*settings.ISBNHeader) == 0 {
		return nil
	}

	return &settings
}

func SetFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&volumeV, volumeF, "", "header of the column that has the volume number.")
	cmd.PersistentFlags().StringVar(&titleV, titleF, "", "header of the column that has the title.")
	cmd.PersistentFlags().StringVar(&isbnV, isbnF, "", "header of the column that has the ISBN number(required).")
	cmd.PersistentFlags().IntSliceVar(&tableV, tableF, []int{}, "tables to include, zero indexed. skip for all tables.")
}
