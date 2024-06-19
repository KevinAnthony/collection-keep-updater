package libib_test

import (
	"errors"
	"fmt"
	"io"
	native "net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/library/libib"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	Convey("New", t, func() {
		ctx := ctxu.NewContextMock(t)
		cmd := types.NewICommandMock(t)
		httpMock := http.NewClientMock(t)

		settings := types.LibrarySettings{
			Name:        types.LibIBLibrary,
			WantedColID: "id0",
			OtherColIDs: []string{"id1", "id2", "id3"},
			APIKey:      "secret",
		}

		cmd.On("Context").Return(ctx).Once()
		ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(httpMock).Once()

		Convey("should return lib setting", func() {
			So(libib.New(cmd, settings), ShouldNotBeNil)
		})
	})
}

func TestLibIB_GetBooksInCollection(t *testing.T) {
	t.Parallel()

	Convey("GetBooksInCollection", t, func() {
		cmd := types.NewICommandMock(t)
		ctx := ctxu.NewContextMock(t)
		httpMock := http.NewClientMock(t)
		bodyMock := http.NewBodyMock(t)

		settings := types.LibrarySettings{
			Name:        types.LibIBLibrary,
			WantedColID: "id0",
			APIKey:      "secret",
		}

		cmd.On("Context").Return(ctx).Once()
		ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(httpMock).Once()
		getCSVCall := httpMock.On("Do", mock.MatchedBy(matchFunc("id0"))).Maybe()

		client := libib.New(cmd, settings)

		expected := types.ISBNBooks{
			{
				ISBN10: "1626923485",
				ISBN13: "9781626923485",
				Title:  "Miss Kobayashi's Dragon Maid, Vol. 1",
			}, {
				ISBN10: "1626924317",
				ISBN13: "9781626924314",
				Title:  "Miss Kobayashi's Dragon Maid, Vol. 2",
			},
		}

		file, err := openFile(t, "test_fixtures.csv")
		So(err, ShouldBeNil)

		Convey("should return valid collection of books", func() {
			getCSVCall.Once().Return(file, nil)

			actual, err := client.GetBooksInCollection(ctx)

			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})
		Convey("should return error when", func() {
			Convey("do request returns error", func() {
				getCSVCall.Once().Return(nil, errors.New("http do error"))

				actual, err := client.GetBooksInCollection(ctx)

				So(err, ShouldBeError, "http do error")
				So(actual, ShouldBeNil)
			})
			Convey("csv unmarshal fails", func() {
				getCSVCall.Once().Return(bodyMock, nil)
				bodyMock.On("Read", mock.Anything).Once().Return(0, errors.New("everybody body mock"))

				actual, err := client.GetBooksInCollection(ctx)

				So(err, ShouldBeError, "everybody body mock")
				So(actual, ShouldBeNil)
			})
		})
	})
}

func TestLibIB_SaveWanted(t *testing.T) {
	t.Parallel()

	Convey("SaveWanted", t, func() {
		t.Cleanup(func() {
			wd, _ := os.Getwd()
			if _, err := os.Stat(filepath.Join(wd, "wanted.csv")); err == nil {
				_ = os.Remove(filepath.Join(wd, "wanted.csv"))
			}
			if _, err := os.Stat(filepath.Join(wd, "library", "libib", "./wanted.csv")); err == nil {
				_ = os.Remove(filepath.Join(wd, "library", "libib", "wanted.csv"))
			}
		})

		cmd := types.NewICommandMock(t)
		ctx := ctxu.NewContextMock(t)
		httpMock := http.NewClientMock(t)

		settings := types.LibrarySettings{
			Name:        types.LibIBLibrary,
			WantedColID: "id0",
			APIKey:      "secret",
		}

		cmd.On("Context").Return(ctx).Once()
		ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(httpMock).Once()

		expected := types.ISBNBooks{
			{
				ISBN10: "1626923485",
				ISBN13: "9781626923485",
				Title:  "Miss Kobayashi's Dragon Maid, Vol. 1",
			}, {
				ISBN10: "1626924317",
				ISBN13: "9781626924314",
				Title:  "Miss Kobayashi's Dragon Maid, Vol. 2",
			},
		}

		client := libib.New(cmd, settings)
		Convey("should return valid collection of books", func() {
			err := client.SaveWanted(cmd, expected)

			So(err, ShouldBeNil)
		})
	})
}

func matchFunc(id string) func(req *native.Request) bool {
	return func(req *native.Request) bool {
		s, err := io.ReadAll(req.Body)
		So(err, ShouldBeNil)

		return strings.HasSuffix(string(s), fmt.Sprintf("settings-library-export-id=%s", id))
	}
}

func openFile(t *testing.T, fileName string) (*os.File, error) {
	t.Helper()

	wd, err := os.Getwd()
	So(err, ShouldBeNil)

	if !strings.HasSuffix(wd, "libib") {
		return os.Open(filepath.Join(wd, "library", "libib", fileName))
	}

	return os.Open(filepath.Join(wd, fileName))
}
