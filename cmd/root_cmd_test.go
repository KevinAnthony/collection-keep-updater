package cmd_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/cmd"
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/di"
	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetRootCmd(t *testing.T) {
	t.Parallel()
	Convey("GetRootCmd", t, func() {
		Convey("should return non-nil command", func() {
			So(cmd.GetRootCmd(), ShouldNotBeNil)
		})
	})
}

func TestPreRunE(t *testing.T) {
	t.Parallel()

	Convey("TestPreRunE", t, func() {
		ctx := ctxu.NewContextMock(t)
		command := types.NewICommandMock(t)
		factory := di.NewIDepFactoryMock(t)
		cfgLoader := types.NewIConfigMock(t)

		_ = factory
		command.On("Context").Return(ctx).Times(2)
		ctx.On("Value", ctxu.ContextKey("config_loader_ctx_key")).Return(cfgLoader).Once()
		ctx.On("Value", ctxu.ContextKey("dep_factory_ctx_key")).Return(factory).Once()
		cfgCall := factory.On("Config", command, cfgLoader).Maybe()
		srcCall := factory.On("Sources", command).Maybe()
		libCall := factory.On("Libraries", command).Maybe()

		Convey("should return no errors", func() {
			cfgCall.Return(nil).Once()
			srcCall.Return(nil).Once()
			libCall.Return(nil).Once()

			err := cmd.PreRunE(command)

			So(err, ShouldBeNil)
		})
		Convey("should return error when", func() {
			Convey("config returns an error", func() {
				cfgCall.Return(errors.New("cfg error")).Once()

				err := cmd.PreRunE(command)

				So(err, ShouldBeError, "cfg error")
			})
			Convey("sources returns an error", func() {
				cfgCall.Return(nil).Once()
				srcCall.Return(errors.New("source error")).Once()

				err := cmd.PreRunE(command)

				So(err, ShouldBeError, "source error")
			})
			Convey("libraries returns an error", func() {
				cfgCall.Return(nil).Once()
				srcCall.Return(nil).Once()
				libCall.Return(errors.New("lib error")).Once()

				err := cmd.PreRunE(command)

				So(err, ShouldBeError, "lib error")
			})
		})
	})
}
