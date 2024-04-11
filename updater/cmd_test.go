package updater_test

import (
	"bytes"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/updater"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/pflag"
)

func TestGetCmd(t *testing.T) {
	t.Parallel()

	Convey("GetCmd", t, func() {
		Convey("should return valid command", func() {
			cmd := updater.GetCmd()

			So(cmd, ShouldNotBeNil)
			So(cmd.Name(), ShouldResemble, "update")
		})
	})
}

func TestRun(t *testing.T) {
	t.Parallel()

	Convey("Run", t, func() {
		ctx := ctxu.NewContextMock(t)
		cmd := types.NewICommandMock(t)
		updateMock := updater.NewIUpdaterMock(t)
		libMock := types.NewILibraryMock(t)

		series := []types.Series{{Name: "test series"}}
		cfg := types.Config{Series: series}
		libs := map[types.LibraryType]types.ILibrary{types.LibIBLibrary: libMock}
		books := types.ISBNBooks{{Title: "test book", ISBN13: "1234567890abc"}, {Title: "test book 2", ISBN13: "1234567890abd"}}
		wanted := types.ISBNBooks{{Title: "test book", ISBN13: "1234567890abc"}}
		flags := &pflag.FlagSet{}

		getCtxCall := cmd.On("Context").Return(ctx).Maybe()
		getCfgCall := ctx.On("Value", ctxu.ContextKey("config_ctx_key")).Maybe()
		getLibCall := ctx.On("Value", ctxu.ContextKey("libraries_ctx_key")).Maybe()
		getUpdaterCall := ctx.On("Value", ctxu.ContextKey("updater_ctx_key")).Return(updateMock).Maybe()
		getBooksCall := updateMock.On("GetAllAvailableBooks", cmd, series).Maybe()
		updateCall := updateMock.On("UpdateLibrary", ctx, libMock, books).Maybe()
		flagCall := cmd.On("PersistentFlags").Maybe()
		saveCall := libMock.On("SaveWanted", wanted).Maybe()

		Convey("should get books and list wanted", func() {
			Convey("and print to screen when print flag set", func() {
				buff := bytes.NewBufferString("")
				flags.Bool("print-config", true, "")
				flags.Bool("write-config", false, "")

				getUpdaterCall.Once()
				getCtxCall.Times(4)
				getCfgCall.Once().Return(cfg, nil)
				getLibCall.Once().Return(libs, nil)
				getBooksCall.Once().Return(books, nil)
				updateCall.Once().Return(wanted, nil)
				flagCall.Once().Return(flags)
				cmd.On("OutOrStdout").Once().Return(buff)

				err := updater.Run(cmd)

				So(err, ShouldBeNil)
				So(buff.String(), ShouldResemble, `┌───────────┬────────┬─────────┬───────────────┬────────┐
│ TITLE     │ VOLUME │ ISBN 10 │ ISBN 13       │ SOURCE │
├───────────┼────────┼─────────┼───────────────┼────────┤
│ test book │        │         │ 1234567890abc │        │
└───────────┴────────┴─────────┴───────────────┴────────┘
`)
			})
			Convey("and call save when save flag set", func() {
				flags.Bool("print-config", false, "")
				flags.Bool("write-config", true, "")

				getCtxCall.Times(4)
				getUpdaterCall.Once()
				getCfgCall.Once().Return(cfg, nil)
				getLibCall.Once().Return(libs, nil)
				getBooksCall.Once().Return(books, nil)
				updateCall.Once().Return(wanted, nil)
				flagCall.Twice().Return(flags)
				saveCall.Once().Return(nil)

				err := updater.Run(cmd)

				So(err, ShouldBeNil)
			})
			Convey("and wanted was zero length", func() {
				getCtxCall.Times(4)
				getUpdaterCall.Once()
				getCfgCall.Once().Return(cfg, nil)
				getLibCall.Once().Return(libs, nil)
				getBooksCall.Once().Return(books, nil)
				updateCall.Once().Return(types.ISBNBooks{}, nil)

				err := updater.Run(cmd)

				So(err, ShouldBeNil)
			})
		})
		Convey("should return error when", func() {
			Convey("get config returns error", func() {
				getCtxCall.Once()
				getCfgCall.Once().Return(nil, errors.New("get config error"))

				err := updater.Run(cmd)

				So(err, ShouldBeError, "configuration not found in context")
			})
			Convey("get libraries returns error", func() {
				getCtxCall.Twice()
				getCfgCall.Once().Return(cfg, nil)
				getLibCall.Once().Return(nil, errors.New("library error"))

				err := updater.Run(cmd)

				So(err, ShouldBeError, "libraries not found in context")
			})
			Convey("get available books returns an error", func() {
				getCtxCall.Times(3)
				getUpdaterCall.Once()
				getCfgCall.Once().Return(cfg, nil)
				getLibCall.Once().Return(libs, nil)
				getBooksCall.Once().Return(nil, errors.New("get book error"))

				err := updater.Run(cmd)

				So(err, ShouldBeError, "get book error")
			})
			Convey("update lib returns an error", func() {
				getCtxCall.Times(4)
				getUpdaterCall.Once()
				getCfgCall.Once().Return(cfg, nil)
				getLibCall.Once().Return(libs, nil)
				getBooksCall.Once().Return(books, nil)
				updateCall.Once().Return(nil, errors.New("update book error"))

				err := updater.Run(cmd)

				So(err, ShouldBeError, "update book error")
			})
			Convey("save returns an error", func() {
				flags.Bool("print-config", false, "")
				flags.Bool("write-config", true, "")

				getCtxCall.Times(4)
				getUpdaterCall.Once()
				getCfgCall.Once().Return(cfg, nil)
				getLibCall.Once().Return(libs, nil)
				getBooksCall.Once().Return(books, nil)
				updateCall.Once().Return(wanted, nil)
				flagCall.Twice().Return(flags)
				saveCall.Once().Return(errors.New("save books error"))

				err := updater.Run(cmd)

				So(err, ShouldBeError, "save books error")
			})
		})
	})
}
