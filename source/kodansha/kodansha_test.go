package kodansha_test

import (
	native "net/http"
	"os"
	"path"
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
		client := http.NewClientMock(t)
		Convey("should return isource when http client is valid", func() {
			source, err := kodansha.New(client)

			So(source, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return error when http client is nil", func() {
			source, err := kodansha.New(nil)

			So(err, ShouldBeError, "http client is nil")
			So(source, ShouldBeNil)
		})
	})
}

func TestKodansha_GetISBNs(t *testing.T) {
	t.Parallel()

	Convey("GetISBNs", t, func() {
		id := "initial-d-test"
		client := http.NewClientMock(t)
		ctx := ctxu.NewContextMock(t)
		bodyMock := http.NewBodyMock(t)

		expected := types.ISBNBooks{{ISBN13: "9798888770986", Volume: "1", Title: "Initial D Omnibus, Volume 1", Source: "Kodansha"}}

		source, err := kodansha.New(client)
		So(err, ShouldBeNil)

		seriesCall := client.On("Do", mock.MatchedBy(func(req *native.Request) bool {
			return strings.HasSuffix(req.URL.String(), id)
		})).Maybe()

		page1Call := client.On("Do", mock.MatchedBy(func(req *native.Request) bool {
			return strings.HasSuffix(req.URL.String(), "initial-d-omnibus-1")
		})).Maybe()

		Convey("should return valid isbn", func() {
			seriesCall.Return(os.Open(getPath(t, "series.html")))
			page1Call.Return(os.Open(getPath(t, "page_1.html")))

			books, err := source.GetISBNs(ctx, types.Series{ID: id})

			So(err, ShouldBeNil)
			So(books, ShouldResemble, expected)
		})
		Convey("should return nil error when series page does not contain any links to the walk page returns nil", func() {
			seriesCall.Return(os.Open(getPath(t, "series_missing_link.html")))

			books, err := source.GetISBNs(ctx, types.Series{ID: id})

			So(err, ShouldBeNil)
			So(books, ShouldBeEmpty)
		})

		Convey("should return error when", func() {
			Convey("id is invalid and new create fails", func() {
				books, err := source.GetISBNs(ctx, types.Series{ID: string([]byte{0x7f})}) // this triggers failed parse

				So(err, ShouldBeError, `parse "https://kodansha.us/series/\x7f": net/url: invalid control character in URL`)
				So(books, ShouldBeNil)
			})
			Convey("do request on series page fails", func() {
				seriesCall.Return(nil, errors.New("do request error"))

				books, err := source.GetISBNs(ctx, types.Series{ID: id})

				So(err, ShouldBeError, `do request error`)
				So(books, ShouldBeNil)
			})
			Convey("http parse on series page fails", func() {
				bodyMock.On("Read", mock.Anything).Return(0, errors.New("everybody body mock"))
				seriesCall.Return(bodyMock, nil)

				books, err := source.GetISBNs(ctx, types.Series{ID: id})

				So(err, ShouldBeError, "everybody body mock")
				So(books, ShouldBeNil)
			})
			Convey("book page is invalid and create request fails", func() {
				seriesCall.Return(os.Open(getPath(t, "series_bad_url.html")))

				books, err := source.GetISBNs(ctx, types.Series{ID: id})

				So(err, ShouldBeError)
				So(err.Error(), ShouldEndWith, `parse "https://kodansha.us/product/initial-d-omnibus-1\x7f": net/url: invalid control character in URL`)
				So(books, ShouldBeNil)
			})
			Convey("book page do request fails", func() {
				seriesCall.Return(os.Open(getPath(t, "series.html")))
				page1Call.Return(nil, errors.New("page request error"))

				books, err := source.GetISBNs(ctx, types.Series{ID: id})

				So(err, ShouldBeError, "/product/initial-d-omnibus-1: page request error")
				So(books, ShouldBeNil)
			})
			Convey("book page html parse fails", func() {
				seriesCall.Return(os.Open(getPath(t, "series.html")))
				bodyMock.On("Read", mock.Anything).Return(0, errors.New("everybody body mock"))
				page1Call.Return(bodyMock, nil)

				books, err := source.GetISBNs(ctx, types.Series{ID: id})

				So(err, ShouldBeError, "/product/initial-d-omnibus-1: everybody body mock")
				So(books, ShouldBeNil)
			})
		})
		Convey("should partially fail when", func() {
			Convey("get volume fails because", func() {
				expected[0].Volume = ""
				Convey("url does not contain enough dashes", func() {
					seriesCall.Return(os.Open(getPath(t, "series_volume_bad_url.html")))
					page1Call.Return(os.Open(getPath(t, "page_1.html")))

					books, err := source.GetISBNs(ctx, types.Series{ID: id})

					So(err, ShouldBeNil)
					So(books, ShouldResemble, expected)
				})
				Convey("volume is not a number", func() {
					seriesCall.Return(os.Open(getPath(t, "series_volume_bad_number.html")))
					client.On("Do", mock.MatchedBy(func(req *native.Request) bool {
						return strings.HasSuffix(req.URL.String(), "initial-d-omnibus-1x")
					})).Return(os.Open(getPath(t, "page_1.html")))

					books, err := source.GetISBNs(ctx, types.Series{ID: id})

					So(err, ShouldBeNil)
					So(books, ShouldResemble, expected)
				})
			})
		})
	})
}

func getPath(t *testing.T, fileName string) string {
	t.Helper()

	wd, err := os.Getwd()
	So(err, ShouldBeNil)

	if !strings.HasSuffix(wd, path.Join("source", "kodansha")) {
		return path.Join(wd, "source", "kodansha", "test_fixtures", fileName)
	}

	return path.Join(wd, "test_fixtures", fileName)
}
