package cmd_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/cmd"
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestGetRootCmd(t *testing.T) {
	t.Parallel()
	Convey("GetRootCmd", t, func() {
		Convey("should return non-nil command", func() {
			So(cmd.GetRootCmd(), ShouldNotBeNil)
		})
	})
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	Convey("LoadConfig", t, func() {
		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)
		cfgMock := types.NewIConfigMock(t)
		source := types.NewISouceMock(t)

		settingMap := map[string]any{
			"delay_between":   "100ms",
			"maximum_backlog": 2,
		}
		seriesBlob := []any{
			map[string]any{
				"id":              "one-piece",
				"key":             "one-piece",
				"name":            "One Piece",
				"source":          "viz",
				"source_settings": settingMap,
				"isbn_blacklist":  []any{"one", "two", "five"},
			},
		}
		libraryBlob := []any{
			map[string]any{
				"api_key":              "secret",
				"other_collection_ids": []any{"id1", "id2", "id3"},
				"wanted_collection_id": "id0",
				"type":                 "libib",
			},
		}
		vizSrc, err := viz.New(http.NewClientMock(t))
		So(err, ShouldBeNil)

		sourceSetting := vizSrc.SourceSettingFromConfig(settingMap)

		expectedCfg := types.Config{
			Series: []types.Series{
				{
					Name:           "One Piece",
					ID:             "one-piece",
					Key:            "one-piece",
					Source:         "viz",
					SourceSettings: sourceSetting,
					ISBNBlacklist:  []string{"one", "two", "five"},
				},
			},
			Libraries: []types.LibrarySettings{
				{
					APIKey:      "secret",
					WantedColID: "id0",
					OtherColIDs: []string{"id1", "id2", "id3"},
					Name:        "libib",
				},
			},
		}

		cfgMock.On("AddConfigPath", "$HOME/.config/noside/").Once()
		cfgMock.On("AddConfigPath", ".").Once()
		cfgMock.On("SetConfigType", "yaml").Once()
		cfgMock.On("SetConfigName", "config").Once()
		cfgMock.On("AutomaticEnv").Once()
		readCall := cfgMock.On("ReadInConfig").Maybe()
		getSeriesCall := cfgMock.On("Get", "series").Maybe()
		getLibCall := cfgMock.On("Get", "libraries").Maybe()
		cmdCall := cmdMock.On("Context").Maybe()
		ctxCall := ctx.On("Value", ctxu.ContextKey("sources_ctx_key")).Maybe()
		sourceCall := source.On("SourceSettingFromConfig", settingMap).Maybe()
		cmdMock.On("SetContext", mock.MatchedBy(matchFunc(cfgMock, expectedCfg))).Maybe()

		sources := map[types.SourceType]types.ISource{types.VizSource: source}

		Convey("should load cfg and set it in context", func() {
			readCall.Return(nil).Once()
			getSeriesCall.Return(seriesBlob).Once()
			cmdCall.Return(ctx).Twice()
			ctxCall.Return(sources).Once()
			sourceCall.Return(sourceSetting).Once()

			getLibCall.Return(libraryBlob).Once()

			err := cmd.LoadConfig(cmdMock, cfgMock)

			So(err, ShouldBeNil)
		})
		Convey("should return error when", func() {
			Convey("read config returns error", func() {
				readCall.Return(errors.New("error here"))

				err := cmd.LoadConfig(cmdMock, cfgMock)

				So(err, ShouldBeError, "error here")
			})
			Convey("get series returns error", func() {
				readCall.Return(nil)
				getSeriesCall.Return([]any{"hello"})

				err := cmd.LoadConfig(cmdMock, cfgMock)

				So(err, ShouldBeError, "data is not a series")
			})
			Convey("get library returns error", func() {
				readCall.Return(nil)
				getSeriesCall.Return(seriesBlob)
				cmdCall.Return(ctx).Once()
				ctxCall.Return(sources).Once()
				sourceCall.Return(sourceSetting).Once()

				getLibCall.Return([]any{"hello"}).Once()

				err := cmd.LoadConfig(cmdMock, cfgMock)

				So(err, ShouldBeError, "data is not a library")
			})
		})
	})
}

func TestLoadDI(t *testing.T) {
	t.Parallel()
	Convey("LoadDI", t, func() {})
}

func matchFunc(expectedIConfig types.IConfig, expectedConfig types.Config) func(ctx context.Context) bool {
	return func(ctx context.Context) bool {
		actualIConfig := ctx.Value(ctxu.ContextKey("i_config_ctx_key"))
		actualConfig := ctx.Value(ctxu.ContextKey("config_ctx_key"))

		return reflect.DeepEqual(actualIConfig, expectedIConfig) &&
			reflect.DeepEqual(actualConfig, expectedConfig)
	}
}
