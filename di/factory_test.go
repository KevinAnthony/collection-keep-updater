package di_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/di"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestDepFactory_Config(t *testing.T) {
	t.Parallel()

	Convey("Config", t, func() {
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

		factory := di.NewDepFactory()

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
		cmdMock.On("SetContext", mock.MatchedBy(matchFunc(expectedCfg))).Maybe()

		sources := map[types.SourceType]types.ISource{types.VizSource: source}

		Convey("should load cfg and set it in context", func() {
			readCall.Return(nil).Once()
			getSeriesCall.Return(seriesBlob).Once()
			cmdCall.Return(ctx).Twice()
			ctxCall.Return(sources).Once()
			sourceCall.Return(sourceSetting).Once()
			getLibCall.Return(libraryBlob).Once()

			err := factory.Config(cmdMock, cfgMock)

			So(err, ShouldBeNil)
		})
		Convey("should return error when", func() {
			Convey("read config returns error", func() {
				readCall.Return(errors.New("error here"))

				err := factory.Config(cmdMock, cfgMock)

				So(err, ShouldBeError, "error here")
			})
			Convey("get series returns error", func() {
				readCall.Return(nil)
				getSeriesCall.Return([]any{"hello"})

				err := factory.Config(cmdMock, cfgMock)

				So(err, ShouldBeError, "data is not a series")
			})
			Convey("get library returns error", func() {
				readCall.Return(nil)
				getSeriesCall.Return(seriesBlob)
				cmdCall.Return(ctx).Once()
				ctxCall.Return(sources).Once()
				sourceCall.Return(sourceSetting).Once()

				getLibCall.Return([]any{"hello"}).Once()

				err := factory.Config(cmdMock, cfgMock)

				So(err, ShouldBeError, "data is not a library")
			})
		})
	})
}

func TestLoadDI(t *testing.T) {
	t.Parallel()
	Convey("LoadSources", t, func() {
		factory := di.NewDepFactory()

		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)
		clientMock := http.NewClientMock(t)
		wikiMock := wikipedia.NewTableGetterMock(t)

		getCall := cmdMock.On("Context").Maybe()
		setCall := cmdMock.On("SetContext", mock.Anything).Maybe()

		Convey("should set source in context", func() {
			getCall.Once().Return(ctx)
			setCall.Once().Return()

			err := factory.Sources(cmdMock, clientMock, wikiMock)

			So(err, ShouldBeNil)
		})

		Convey("should return error when", func() {
			Convey("client mock is nil", func() {
				err := factory.Sources(cmdMock, nil, wikiMock)

				So(err, ShouldBeError, "http client is nil")
			})
			Convey("wiki mock is nil", func() {
				err := factory.Sources(cmdMock, clientMock, nil)

				So(err, ShouldBeError, "wikipedia table getter is nil")
			})
		})
	})
}

func TestGetDIFactory(t *testing.T) {
	t.Parallel()

	Convey("GetDIFactory", t, func() {
		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)
		diMock := di.NewIDepFactoryMock(t)

		getCtxCall := cmdMock.On("Context").Once()

		ctxCall := ctx.On("Value", ctxu.ContextKey("dep_factory_ctx_key")).Maybe()

		Convey("should get this cfg from the context of the command", func() {
			getCtxCall.Return(ctx)
			ctxCall.Once().Return(diMock)

			actual := di.GetDIFactory(cmdMock)

			So(actual, ShouldResemble, diMock)
		})
		Convey("should return new viper when read is not in config", func() {
			cmdMock.On("SetContext", mock.Anything)

			getCtxCall.Return(context.Background())

			actual := di.GetDIFactory(cmdMock)

			So(actual, ShouldHaveSameTypeAs, di.NewDepFactory())
		})
	})
}

func matchFunc(expectedConfig types.Config) func(ctx context.Context) bool {
	return func(ctx context.Context) bool {
		actualConfig := ctx.Value(ctxu.ContextKey("config_ctx_key"))

		return reflect.DeepEqual(actualConfig, expectedConfig)
	}
}
