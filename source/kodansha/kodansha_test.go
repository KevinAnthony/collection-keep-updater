package kodansha_test

import (
	native "net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/source/kodansha"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	Convey("New", t, func() {
		cmd := types.NewICommandMock(t)
		ctx := ctxu.NewContextMock(t)
		client := http.NewClientMock(t)

		cmd.On("Context").Return(ctx)
		ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(client)

		Convey("should return isource", func() {
			source := kodansha.New(cmd)

			So(source, ShouldNotBeNil)
		})
	})
}

func TestKodansha_GetISBNs(t *testing.T) {
	t.Parallel()

	Convey("GetISBNs", t, func() {
		id := "initial-d-test"
		client := http.NewClientMock(t)
		cmd := types.NewICommandMock(t)
		ctx := ctxu.NewContextMock(t)
		bodyMock := http.NewBodyMock(t)

		series := types.Series{ID: id}
		expected := types.ISBNBooks{{ISBN13: "9798888770986", Volume: "1", Title: "Initial D Omnibus, Volume 1", Source: "Kodansha"}}

		cmd.On("Context").Return(ctx)
		ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(client)

		source := kodansha.New(cmd)

		seriesCall := client.On("Do", mock.MatchedBy(matchFunc(id))).Maybe()

		page1Call := client.On("Do", mock.MatchedBy(matchFunc("initial-d-omnibus-1"))).Maybe()

		Convey("should return valid isbn", func() {
			seriesCall.Return(openFile(t, "series.html")).Once()
			page1Call.Return(openFile(t, "page_1.html")).Once()

			books, err := source.GetISBNs(ctx, series)

			So(err, ShouldBeNil)
			So(books, ShouldResemble, expected)
		})
		Convey("should return nil error when series page does not contain any links to the walk page returns nil", func() {
			seriesCall.Return(openFile(t, "series_missing_link.html")).Once()

			books, err := source.GetISBNs(ctx, series)

			So(err, ShouldBeNil)
			So(books, ShouldBeEmpty)
		})

		Convey("should return error when", func() {
			Convey("for series", func() {
				Convey("id is invalid and new create fails", func() {
					books, err := source.GetISBNs(ctx, types.Series{ID: string([]byte{0x7f})}) // this triggers failed parse

					So(err, ShouldBeError, `parse "https://kodansha.us/series/\x7f": net/url: invalid control character in URL`)
					So(books, ShouldBeNil)
				})
				Convey("do request fails", func() {
					seriesCall.Return(nil, errors.New("do request error")).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, `do request error`)
					So(books, ShouldBeNil)
				})
				Convey("http parse fails", func() {
					bodyMock.On("Read", mock.Anything).Return(0, errors.New("everybody body mock")).Once()
					seriesCall.Return(bodyMock, nil).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, "everybody body mock")
					So(books, ShouldBeNil)
				})
			})
			Convey("for book page", func() {
				Convey("is invalid and create request fails", func() {
					seriesCall.Return(openFile(t, "series_bad_url.html")).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError)
					So(err.Error(), ShouldEndWith, `parse "https://kodansha.us/product/initial-d-omnibus-1\x7f": net/url: invalid control character in URL`)
					So(books, ShouldBeNil)
				})
				Convey("do request fails", func() {
					seriesCall.Return(openFile(t, "series.html")).Once()
					page1Call.Return(nil, errors.New("page request error")).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, "/product/initial-d-omnibus-1: page request error")
					So(books, ShouldBeNil)
				})
				Convey("html parse fails", func() {
					seriesCall.Return(openFile(t, "series.html")).Once()
					bodyMock.On("Read", mock.Anything).Return(0, errors.New("everybody body mock")).Once()
					page1Call.Return(bodyMock, nil)

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, "/product/initial-d-omnibus-1: everybody body mock")
					So(books, ShouldBeNil)
				})
			})
		})
		Convey("should partially fail when", func() {
			Convey("get volume fails because", func() {
				expected[0].Volume = ""
				Convey("url does not contain enough dashes", func() {
					seriesCall.Return(openFile(t, "series_volume_bad_url.html")).Once()
					page1Call.Return(openFile(t, "page_1.html")).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeNil)
					So(books, ShouldResemble, expected)
				})
				Convey("volume is not a number", func() {
					seriesCall.Return(openFile(t, "series_volume_bad_number.html")).Once()
					client.On("Do", mock.MatchedBy(matchFunc("initial-d-omnibus-1x"))).
						Return(openFile(t, "page_1.html")).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeNil)
					So(books, ShouldResemble, expected)
				})
			})
		})
	})
}

func matchFunc(id string) func(req *native.Request) bool {
	return func(req *native.Request) bool {
		s := req.URL.String()
		return strings.HasSuffix(s, id)
	}
}

func openFile(t *testing.T, fileName string) (*os.File, error) {
	t.Helper()

	wd, err := os.Getwd()
	So(err, ShouldBeNil)

	if !strings.HasSuffix(wd, filepath.Join("source", "kodansha")) {
		return os.Open(filepath.Join(wd, "source", "kodansha", "test_fixtures", fileName))
	}

	return os.Open(filepath.Join(wd, "test_fixtures", fileName))
}
