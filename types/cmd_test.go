package types_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
)

func TestCmdArgs(t *testing.T) {
	t.Parallel()

	Convey("CmdArgs", t, func() {
		cmd := &cobra.Command{}

		Convey("should call f with cmd, and return error if func returns error", func() {
			f := func(iCommand types.ICommand, args []string) error {
				So(cmd, ShouldEqual, iCommand)

				return errors.New("this was called")
			}
			fn := types.CmdArgs(f)

			err := fn(cmd, nil)

			So(err, ShouldBeError, "this was called")
		})
	})
}

func TestCmdPersistentPreRunE(t *testing.T) {
	t.Parallel()

	Convey("CmdPersistentPreRunE", t, func() {
		cmd := &cobra.Command{}

		Convey("should call f with cmd, and return error if func returns error", func() {
			f := func(iCommand types.ICommand, args []string) error {
				So(cmd, ShouldEqual, iCommand)

				return errors.New("this was called")
			}
			fn := types.CmdPersistentPreRunE(f)

			err := fn(cmd, nil)

			So(err, ShouldBeError, "this was called")
		})
	})
}

func TestCmdRunE(t *testing.T) {
	Convey("CmdRunE", t, func() {
		cmd := &cobra.Command{}

		Convey("should call f with cmd, and return error if func returns error", func() {
			f := func(iCommand types.ICommand, args []string) error {
				So(cmd, ShouldEqual, iCommand)

				return errors.New("this was called")
			}
			fn := types.CmdRunE(f)

			err := fn(cmd, nil)

			So(err, ShouldBeError, "this was called")
		})
	})
}
