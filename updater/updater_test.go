package updater_test

import (
	"errors"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/updater"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewUpdater(t *testing.T) {
	t.Parallel()

	Convey("NewUpdater", t, func() {
		update := updater.NewUpdater()

		So(update, ShouldNotBeNil)
	})
}

func TestUpdater_GetAllAvailableBooks(t *testing.T) {
	t.Parallel()

	Convey("GetAllAvailableBooks", t, func() {
		ctx := ctxu.NewContextMock(t)
		cmdMock := types.NewICommandMock(t)
		sourceMock := types.NewISouceMock(t)
		update := updater.NewUpdater()

		expected := types.ISBNBooks{{Title: "test vol 1", ISBN13: "test_1_isbn"}}
		series := []types.Series{{ID: "test-id", Source: types.VizSource, Name: "test name"}}
		sources := map[types.SourceType]types.ISource{types.VizSource: sourceMock}

		getCtxCall := cmdMock.On("Context").Maybe().Return(ctx)
		getSourceCall := ctx.On("Value", ctxu.ContextKey("sources_ctx_key")).Maybe()
		getBooksCall := sourceMock.On("GetISBNs", ctx, series[0]).Maybe()

		Convey("should return available books", func() {
			getCtxCall.Twice()
			getSourceCall.Once().Return(sources)
			getBooksCall.Return(expected, nil)

			actual, err := update.GetAllAvailableBooks(cmdMock, series)

			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})
		Convey("should skip series when id is empty", func() {
			series[0].ID = ""

			actual, err := update.GetAllAvailableBooks(cmdMock, series)

			So(err, ShouldBeNil)
			So(actual, ShouldBeEmpty)
		})
		Convey("should return error when", func() {
			Convey("get source returns an error", func() {
				series[0].Source = types.YenSource
				getCtxCall.Once()
				getSourceCall.Once().Return(sources)

				actual, err := update.GetAllAvailableBooks(cmdMock, series)

				So(err, ShouldBeError, "yen is unknown: source type yen not found in source map")
				So(actual, ShouldBeEmpty)
			})
			Convey("get isbn returns an error", func() {
				getCtxCall.Twice()
				getSourceCall.Once().Return(sources)
				getBooksCall.Return(nil, errors.New("get book error"))

				actual, err := update.GetAllAvailableBooks(cmdMock, series)

				So(err, ShouldBeError, "test name: get book error")
				So(actual, ShouldBeEmpty)
			})
		})
	})
}

func TestUpdater_UpdateLibrary(t *testing.T) {
	t.Parallel()

	Convey("UpdateLibrary", t, func() {
		ctx := ctxu.NewContextMock(t)
		libMock := types.NewILibraryMock(t)
		wanted := types.ISBNBooks{
			{Title: "test vol 1", ISBN13: "test_1_isbn"},
			{Title: "test vol 2", ISBN13: "test_2_isbn"},
			{Title: "test vol 3", ISBN13: "test_3_isbn"},
			{Title: "test vol 4", ISBN13: "test_4_isbn"},
		}
		inCollection := wanted[1:]

		update := updater.NewUpdater()

		getBooksCall := libMock.On("GetBooksInCollection", ctx).Maybe()

		Convey("should return a diff-ed ", func() {
			getBooksCall.Once().Return(inCollection, nil)

			actual, err := update.UpdateLibrary(ctx, libMock, wanted)

			So(err, ShouldBeNil)
			So(actual, ShouldResemble, wanted[:1])
		})
		Convey("should return error when", func() {
			Convey("get books return an error", func() {
				getBooksCall.Once().Return(nil, errors.New("get books error"))

				actual, err := update.UpdateLibrary(ctx, libMock, wanted)

				So(err, ShouldBeError, "get books error")
				So(actual, ShouldBeNil)
			})
		})
	})
}
