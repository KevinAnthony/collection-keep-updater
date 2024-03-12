package ctxu_test

import (
	"context"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/smartystreets/goconvey/convey"
)

func TestSetConfig(t *testing.T) {
	convey.Convey("SetConfig", t, func() {
		cfg := types.Config{}
		ctx := context.Background()
		expectedCtx := context.WithValue(ctx, ctxu.ContextKey("config_ctx_key"), cfg)
		cmdMock := types.NewICommandMock(t)

		convey.Convey("should set the cfg into the context of the command", func() {
			cmdMock.On("Context").Once().Return(ctx)
			cmdMock.On("SetContext", expectedCtx).Once()

			ctxu.SetConfig(cmdMock, cfg)

			cmdMock.AssertExpectations(t)
		})
	})
}

func TestGetConfig(t *testing.T) {
	convey.Convey("GetConfig", t, func() {
		cfg := types.Config{}
		ctx := context.WithValue(context.Background(), ctxu.ContextKey("config_ctx_key"), cfg)
		cmdMock := types.NewICommandMock(t)

		convey.Convey("should get this cfg from the context of the command", func() {
			cmdMock.On("Context").Once().Return(ctx)

			actual, err := ctxu.GetConfig(cmdMock)

			convey.So(actual, convey.ShouldResemble, cfg)
			convey.So(err, convey.ShouldBeNil)

			cmdMock.AssertExpectations(t)
		})
		convey.Convey("should return error when cfg not in context", func() {
			cmdMock.On("Context").Once().Return(context.Background())

			_, err := ctxu.GetConfig(cmdMock)

			convey.So(err, convey.ShouldBeError, "configuration not found in context")

			cmdMock.AssertExpectations(t)
		})
	})
}
