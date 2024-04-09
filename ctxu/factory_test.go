package ctxu_test

import (
	"context"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
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
