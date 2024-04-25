package di_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/di"
	"github.com/kevinanthony/collection-keep-updater/library/libib"
	"github.com/kevinanthony/collection-keep-updater/source/kodansha"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/source/yen"
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
		cmd := types.NewICommandMock(t)
		client := http.NewClientMock(t)
		wikiMock := wikipedia.NewTableGetterMock(t)

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

		cmd.On("Context").Return(ctx)
		ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(client)

		vizSrc := viz.New(cmd)

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

		ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(client).Maybe()
		ctx.On("Value", ctxu.ContextKey("wiki_getter_ctx_key")).Return(wikiMock).Maybe()

		sources := map[types.SourceType]types.ISource{
			types.WikipediaSource: wikipedia.New(cmd),
			types.VizSource:       viz.New(cmd),
			types.YenSource:       yen.New(cmd),
			types.Kodansha:        kodansha.New(cmd),
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
		getSourceCall := ctx.On("Value", ctxu.ContextKey("sources_ctx_key")).Maybe()
		getConfigCall := ctx.On("Value", ctxu.ContextKey("config_ctx_key")).Maybe()
		// sourceCall := source.On("SourceSettingFromConfig", settingMap).Maybe()
		cmdMock.On("SetContext", mock.MatchedBy(matchFunc("config_ctx_key", expectedCfg))).Maybe()
		cmdMock.On("SetContext", mock.MatchedBy(matchFunc("sources_ctx_key", sources))).Maybe()

		Convey("should load cfg and set it in context", func() {
			readCall.Return(nil).Once()
			getSeriesCall.Return(seriesBlob).Once()
			cmdCall.Return(ctx).Times(8)
			getSourceCall.Return(sources).Once()
			getConfigCall.Return(cfgMock).Once()
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
				cmdCall.Return(ctx).Times(6)
				getSeriesCall.Return([]any{"hello"})

				err := factory.Config(cmdMock, cfgMock)

				So(err, ShouldBeError, "data is not a series")
			})
			Convey("get library returns error", func() {
				readCall.Return(nil)
				getSeriesCall.Return(seriesBlob)
				cmdCall.Return(ctx).Times(7)
				getSourceCall.Return(sources).Once()
				getConfigCall.Return(cfgMock).Once()
				getLibCall.Return([]any{"hello"}).Once()

				err := factory.Config(cmdMock, cfgMock)

				So(err, ShouldBeError, "data is not a library")
			})
		})
	})
}

func TestDepFactory_Libraries(t *testing.T) {
	t.Parallel()

	Convey("Libraries", t, func() {
		factory := di.NewDepFactory()

		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)
		clientMock := http.NewClientMock(t)
		getCall := cmdMock.On("Context").Return(ctx).Maybe()
		ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(clientMock).Maybe()

		settings := types.LibrarySettings{
			Name:        types.LibIBLibrary,
			WantedColID: "id0",
			OtherColIDs: []string{"id1", "id2", "id3"},
			APIKey:      "secret",
		}
		libSettings := map[types.LibraryType]types.ILibrary{types.LibIBLibrary: libib.New(cmdMock, settings)}

		setCall := cmdMock.On("SetContext", mock.MatchedBy(matchFunc("libraries_ctx_key", libSettings))).Maybe()

		getConfigCall := ctx.On("Value", ctxu.ContextKey("config_ctx_key")).Maybe()

		cfg := types.Config{
			Libraries: []types.LibrarySettings{
				{
					APIKey:      "secret",
					WantedColID: "id0",
					OtherColIDs: []string{"id1", "id2", "id3"},
					Name:        "libib",
				},
			},
		}

		Convey("should set source in context", func() {
			getCall.Times(3)
			setCall.Once().Return()
			getConfigCall.Return(cfg, nil).Once()

			err := factory.Libraries(cmdMock)

			So(err, ShouldBeNil)
		})
		Convey("should return error when", func() {
			Convey("get config return error", func() {
				getConfigCall.Return(nil, errors.New("get config error")).Once()

				err := factory.Libraries(cmdMock)

				So(err, ShouldBeError, "configuration not found in context")
			})
		})
	})
}

func matchFunc(key string, expected any) func(ctx context.Context) bool {
	return func(ctx context.Context) bool {
		actual := ctx.Value(ctxu.ContextKey(key))

		return reflect.DeepEqual(actual, expected)
	}
}
