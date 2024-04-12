package utils

import (
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/pkg/errors"
)

func NewLibraryFromMap(_ types.ICommand, data any) (types.LibrarySettings, error) {
	values, ok := data.(map[string]any)
	if !ok {
		return types.LibrarySettings{}, errors.New("data is not a library")
	}
	return types.LibrarySettings{
		Name:        types.LibraryType(Get[string](values, "type")),
		WantedColID: Get[string](values, "wanted_collection_id"),
		OtherColIDs: GetArray[string](values, "other_collection_ids"),
		APIKey:      Get[string](values, "api_key"),
	}, nil
}

func NewSeriesFromMap(cmd types.ICommand, data any) (types.Series, error) {
	values, ok := data.(map[string]any)
	if !ok {
		return types.Series{}, errors.New("data is not a series")
	}

	series := types.Series{
		Name:          Get[string](values, "name"),
		ID:            Get[string](values, "id"),
		Key:           Get[string](values, "key"),
		Source:        types.SourceType(Get[string](values, "source")),
		ISBNBlacklist: GetArray[string](values, "isbn_blacklist"),
	}

	sourceValue := Get[map[string]any](values, "source_settings")
	series.SourceSettings = getSetting(cmd, series.Source, sourceValue)

	return series, nil
}

func getSetting(cmd types.ICommand, key types.SourceType, data map[string]any) types.ISourceSettings {
	if len(data) == 0 {
		return nil
	}

	sourceSetting, err := ctxu.GetSourceSetting(cmd, key)
	if err != nil {
		return nil
	}

	return sourceSetting.SourceSettingFromConfig(data)
}
