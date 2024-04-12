package utils_test

import (
	"testing"

	"github.com/kevinanthony/gorps/v2/http"

	"github.com/kevinanthony/collection-keep-updater/ctxu"

	"github.com/kevinanthony/collection-keep-updater/source/viz"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewLibraryFromMap(t *testing.T) {
	t.Parallel()

	Convey("NewLibraryFromMap", t, func() {
		libraryMap := map[string]any{
			"api_key":              "secret",
			"other_collection_ids": []any{"id1", "id2", "id3"},
			"wanted_collection_id": "id0",
			"type":                 "libib",
		}
		expected := types.LibrarySettings{
			Name:        "libib",
			WantedColID: "id0",
			OtherColIDs: []string{"id1", "id2", "id3"},
			APIKey:      "secret",
		}
		Convey("should return new library settings from map", func() {
			settings, err := utils.NewLibraryFromMap(nil, libraryMap)

			So(err, ShouldBeNil)
			So(settings, ShouldResemble, expected)
		})
		Convey("should return error when", func() {
			Convey("map is not a string/any", func() {
				settings, err := utils.NewLibraryFromMap(nil, []string{"omg", "hi"})

				So(err, ShouldBeError, "data is not a library")
				So(settings, ShouldBeNil)
			})
		})
	})
}

func TestNewSeriesFromMap(t *testing.T) {
	t.Parallel()

	Convey("NewSeriesFromMap", t, func() {
		cmd := types.NewICommandMock(t)
		ctx := ctxu.NewContextMock(t)
		client := http.NewClientMock(t)
		source := types.NewISouceMock(t)

		cmd.On("Context").Return(ctx).Maybe()
		ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(client)

		sources := map[types.SourceType]types.ISource{types.VizSource: source}
		settingMap := map[string]any{
			"delay_between":   "100ms",
			"maximum_backlog": 2,
		}
		seriesMap := map[string]any{
			"id":              "one-piece",
			"key":             "one-piece",
			"name":            "One Piece",
			"source":          "viz",
			"source_settings": settingMap,
			"isbn_blacklist":  []any{"one", "two", "five"},
		}
		vizSrc := viz.New(cmd)
		sourceSetting := vizSrc.SourceSettingFromConfig(settingMap)
		expected := types.Series{
			Name:           "One Piece",
			ID:             "one-piece",
			Key:            "one-piece",
			Source:         "viz",
			SourceSettings: sourceSetting,
			ISBNBlacklist:  []string{"one", "two", "five"},
		}

		getSourceCall := ctx.On("Value", ctxu.ContextKey("sources_ctx_key")).Maybe()
		sourceCall := source.On("SourceSettingFromConfig", settingMap).Maybe()

		Convey("should return series settings", func() {
			getSourceCall.Once().Return(sources)
			sourceCall.Return(sourceSetting).Once()

			actual, err := utils.NewSeriesFromMap(cmd, seriesMap)

			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})
		Convey("should skip source settings when", func() {
			expected.SourceSettings = nil
			Convey("source_settings is nil", func() {
				seriesMap["source_settings"] = nil

				actual, err := utils.NewSeriesFromMap(cmd, seriesMap)

				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
			Convey("get source settings returns nil", func() {
				seriesMap["source"] = "invalid"
				expected.Source = "invalid"

				getSourceCall.Once().Return(sources)

				actual, err := utils.NewSeriesFromMap(cmd, seriesMap)

				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
		})
		Convey("should return error when", func() {
			Convey("map is not string/any", func() {
				actual, err := utils.NewSeriesFromMap(cmd, []string{"omg", "hai"})

				So(err, ShouldBeError, "data is not a series")
				So(actual, ShouldResemble, types.Series{})
			})
		})
	})
}
