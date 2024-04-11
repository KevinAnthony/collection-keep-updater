package di_test

import (
	"context"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/di"
	"github.com/kevinanthony/collection-keep-updater/types"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

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

func TestNewDepFactory(t *testing.T) {
	t.Parallel()

	Convey("NewDepFactory", t, func() {
		So(di.NewDepFactory(), ShouldNotBeNil)
	})
}
