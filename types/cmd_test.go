package types_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
)

func TestCmdRunArgs(t *testing.T) {
	t.Parallel()

	Convey("CmdRunArgs", t, func() {
		cmd := &cobra.Command{}

		Convey("should call f with cmd, and return error if func returns error", func() {
			f := func(iCommand types.ICommand, args []string) error {
				So(cmd, ShouldEqual, iCommand)

				return errors.New("this was called")
			}
			fn := types.CmdRunArgs(f)

			err := fn(cmd, nil)

			So(err, ShouldBeError, "this was called")
		})
	})
}

func TestCmdRun(t *testing.T) {
	Convey("CmdRun", t, func() {
		cmd := &cobra.Command{}

		Convey("should call f with cmd, and return error if func returns error", func() {
			f := func(iCommand types.ICommand) error {
				So(cmd, ShouldEqual, iCommand)

				return errors.New("this was called")
			}
			fn := types.CmdRun(f)

			err := fn(cmd, nil)

			So(err, ShouldBeError, "this was called")
		})
	})
}
