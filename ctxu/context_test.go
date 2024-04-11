package ctxu_test

import (
	"context"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/library/libib"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestSetConfig(t *testing.T) {
	t.Parallel()

	Convey("SetConfig", t, func() {
		cfg := types.Config{}
		ctx := context.Background()
		cmdMock := types.NewICommandMock(t)

		valueCtx := context.WithValue(ctx, ctxu.ContextKey("config_ctx_key"), cfg)

		Convey("should set the cfg into the context of the command", func() {
			cmdMock.On("Context").Once().Return(ctx)
			cmdMock.On("SetContext", valueCtx).Once()

			ctxu.SetConfig(cmdMock, cfg)
		})
	})
}

func TestGetConfig(t *testing.T) {
	t.Parallel()

	Convey("GetConfig", t, func() {
		cfg := types.Config{}
		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)

		cmdCall := cmdMock.On("Context").Once()
		ctxCall := ctx.On("Value", ctxu.ContextKey("config_ctx_key")).Maybe()

		Convey("should get this cfg from the context of the command", func() {
			cmdCall.Return(ctx)
			ctxCall.Once().Return(cfg)

			actual, err := ctxu.GetConfig(cmdMock)

			So(actual, ShouldResemble, cfg)
			So(err, ShouldBeNil)
		})
		Convey("should return error when cfg not in context", func() {
			cmdCall.Return(context.Background())

			_, err := ctxu.GetConfig(cmdMock)

			So(err, ShouldBeError, "configuration not found in context")
		})
	})
}

func TestSetLibraries(t *testing.T) {
	t.Parallel()

	Convey("SetLibraries", t, func() {
		ctx := context.Background()
		httpMock := http.NewClientMock(t)
		settings := types.LibrarySettings{Name: types.LibIBLibrary}

		libSettings := map[types.LibraryType]types.ILibrary{types.LibIBLibrary: libib.New(settings, httpMock)}
		expectedCtx := context.WithValue(context.Background(), ctxu.ContextKey("libraries_ctx_key"), libSettings)
		cmdMock := types.NewICommandMock(t)

		Convey("should set the cfg into the context of the command", func() {
			cmdMock.On("Context").Once().Return(ctx)
			cmdMock.On("SetContext", expectedCtx).Once()

			ctxu.SetLibraries(cmdMock, libSettings)
		})
	})
}

func TestGetLibraries(t *testing.T) {
	t.Parallel()

	Convey("GetLibraries", t, func() {
		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)

		cmdCall := cmdMock.On("Context").Once()
		ctxCall := ctx.On("Value", ctxu.ContextKey("libraries_ctx_key")).Maybe()

		expected := map[types.LibraryType]types.ILibrary{
			types.LibIBLibrary: types.NewILibraryMock(t),
		}
		Convey("should return library when library is in context", func() {
			cmdCall.Return(ctx)
			ctxCall.Return(expected).Once()

			actual, err := ctxu.GetLibraries(cmdMock)

			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})
		Convey("should return error when", func() {
			Convey("library is not in context", func() {
				cmdCall.Return(context.Background())

				_, err := ctxu.GetLibraries(cmdMock)

				So(err, ShouldBeError, "libraries not found in context")
			})
		})
	})
}

func TestSetDI(t *testing.T) {
	t.Parallel()

	Convey("SetDi", t, func() {
		cmdMock := types.NewICommandMock(t)
		ctx := context.Background()

		sources := map[types.SourceType]types.ISource{}

		Convey("should set the cfg into the context of the command", func() {
			cmdMock.On("Context").Once().Return(ctx)
			cmdMock.On("SetContext", mock.Anything).Once()

			ctxu.SetSources(cmdMock, sources)
		})
	})
}

func TestGetSource(t *testing.T) {
	t.Parallel()

	Convey("GetSource", t, func() {
		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)

		cmdCall := cmdMock.On("Context").Once()
		ctxCall := ctx.On("Value", ctxu.ContextKey("sources_ctx_key")).Maybe()

		source := types.NewISouceMock(t)
		sources := map[types.SourceType]types.ISource{types.WikipediaSource: source}

		Convey("should return source when source in context", func() {
			cmdCall.Return(ctx)
			ctxCall.Return(sources).Once()

			actual, err := ctxu.GetSource(cmdMock, types.WikipediaSource)

			So(actual, ShouldResemble, source)
			So(err, ShouldBeNil)
		})
		Convey("should return error when", func() {
			Convey("source type is not in sources", func() {
				cmdCall.Return(context.Background())

				actual, err := ctxu.GetSource(cmdMock, types.WikipediaSource)

				So(actual, ShouldBeNil)
				So(err, ShouldBeError, "sources not found in context")
			})
			Convey("sources is not set in context", func() {
				cmdCall.Return(ctx)
				ctxCall.Return(sources).Once()

				actual, err := ctxu.GetSource(cmdMock, types.VizSource)

				So(actual, ShouldBeNil)
				So(err, ShouldBeError, "source type viz not found in source map")
			})
		})
	})
}

func TestGetSourceSetting(t *testing.T) {
	t.Parallel()

	Convey("GetSourceSetting", t, func() {
		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)

		cmdCall := cmdMock.On("Context").Once()
		ctxCall := ctx.On("Value", ctxu.ContextKey("sources_ctx_key")).Maybe()

		source := types.NewISouceMock(t)
		sources := map[types.SourceType]types.ISource{types.WikipediaSource: source}

		Convey("should return source when source in context", func() {
			cmdCall.Return(ctx)
			ctxCall.Return(sources).Once()

			actual, err := ctxu.GetSourceSetting(cmdMock, types.WikipediaSource)

			So(actual, ShouldResemble, source)
			So(err, ShouldBeNil)
		})
		Convey("should return error when", func() {
			Convey("source type is not in sources", func() {
				cmdCall.Return(context.Background())

				actual, err := ctxu.GetSourceSetting(cmdMock, types.WikipediaSource)

				So(actual, ShouldBeNil)
				So(err, ShouldBeError, "sources not found in context")
			})
			Convey("sources is not set in context", func() {
				cmdCall.Return(ctx)
				ctxCall.Return(sources).Once()

				actual, err := ctxu.GetSourceSetting(cmdMock, types.VizSource)

				So(actual, ShouldBeNil)
				So(err, ShouldBeError, "source type viz not found in source map")
			})
		})
	})
}
