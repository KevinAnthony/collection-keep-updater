package ctxu_test

import (
	"context"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/encoder"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
)

func TestGetConfigReader(t *testing.T) {
	t.Parallel()

	Convey("GetConfigReader", t, func() {
		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)
		cfgMock := types.NewIConfigMock(t)

		getCtxCall := cmdMock.On("Context").Once()

		ctxCall := ctx.On("Value", ctxu.ContextKey("config_loader_ctx_key")).Maybe()

		Convey("should get this cfg from the context of the command", func() {
			getCtxCall.Return(ctx)
			ctxCall.Once().Return(cfgMock)

			actual := ctxu.GetConfigReader(cmdMock)

			So(actual, ShouldResemble, cfgMock)
		})
		Convey("should return new viper when read is not in config", func() {
			cmdMock.On("SetContext", mock.Anything)

			getCtxCall.Return(context.Background())

			actual := ctxu.GetConfigReader(cmdMock)

			So(actual, ShouldHaveSameTypeAs, viper.New())
		})
	})
}

func TestGetHttpClient(t *testing.T) {
	t.Parallel()

	Convey("GetHttpClient", t, func() {
		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)
		httpMock := http.NewClientMock(t)

		getCtxCall := cmdMock.On("Context").Once()

		ctxCall := ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Maybe()

		Convey("should get this cfg from the context of the command", func() {
			getCtxCall.Return(ctx)
			ctxCall.Once().Return(httpMock)

			actual := ctxu.GetHttpClient(cmdMock)

			So(actual, ShouldResemble, httpMock)
		})
		Convey("should return new viper when read is not in config", func() {
			cmdMock.On("SetContext", mock.Anything)

			getCtxCall.Return(context.Background())

			actual := ctxu.GetHttpClient(cmdMock)

			So(actual, ShouldHaveSameTypeAs, http.NewClient(http.NewNativeClient(), encoder.NewFactory()))
		})
	})
}

func TestGetWikiGetter(t *testing.T) {
	t.Parallel()

	Convey("GetWikiGetter", t, func() {
		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)
		wikiMock := wikipedia.NewTableGetterMock(t)

		getCtxCall := cmdMock.On("Context").Once()

		ctxCall := ctx.On("Value", ctxu.ContextKey("wiki_getter_ctx_key")).Maybe()

		Convey("should get this cfg from the context of the command", func() {
			getCtxCall.Return(ctx)
			ctxCall.Once().Return(wikiMock)

			actual := ctxu.GetWikiGetter(cmdMock)

			So(actual, ShouldResemble, wikiMock)
		})
		Convey("should return new viper when read is not in config", func() {
			cmdMock.On("SetContext", mock.Anything)

			getCtxCall.Return(context.Background())

			actual := ctxu.GetWikiGetter(cmdMock)

			So(actual, ShouldNotBeNil)
		})
	})
}
