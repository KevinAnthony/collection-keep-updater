package yen_test

import (
	native "net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/stretchr/testify/mock"

	"github.com/kevinanthony/collection-keep-updater/source/yen"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	t.Parallel()

	Convey("New", t, func() {
		client := http.NewClientMock(t)
		Convey("should return isource when http client is valid", func() {
			source, err := yen.New(client)

			So(source, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return error when http client is nil", func() {
			source, err := yen.New(nil)

			So(err, ShouldBeError, "http client is nil")
			So(source, ShouldBeNil)
		})
	})
}

func TestYen_GetISBNs(t *testing.T) {
	t.Parallel()

	Convey("GetISBNs", t, func() {
		id := "dan-machi"
		client := http.NewClientMock(t)
		ctx := ctxu.NewContextMock(t)
		bodyMock := http.NewBodyMock(t)

		series := types.Series{ID: id}
		expected := types.ISBNBooks{
			types.ISBNBook{ISBN10: "", ISBN13: "9780316339155", Title: " Vol. 1", Volume: "1", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9780316340144", Title: " Vol. 2", Volume: "2", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9780316340151", Title: " Vol. 3", Volume: "3", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9780316340168", Title: " Vol. 4", Volume: "4", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9780316314794", Title: " Vol. 5", Volume: "5", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9780316394161", Title: " Vol. 6", Volume: "6", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9780316394178", Title: " Vol. 7", Volume: "7", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9780316394185", Title: " Vol. 8", Volume: "8", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9780316562645", Title: " Vol. 9", Volume: "9", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9780316442459", Title: " Vol. 10", Volume: "10", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9780316442473", Title: " Vol. 11", Volume: "11", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9781975354787", Title: " Vol. 12", Volume: "12", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9781975328191", Title: " Vol. 13", Volume: "13", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9781975385019", Title: " Vol. 14", Volume: "14", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9781975316105", Title: " Vol. 15", Volume: "15", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9781975333515", Title: " Vol. 16", Volume: "16", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9781975345655", Title: " Vol. 17", Volume: "17", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9781975373917", Title: " Vol. 18", Volume: "18", Source: "Yen Press"},
			types.ISBNBook{ISBN10: "", ISBN13: "9781975393403", Title: " Vol. 19", Volume: "19", Source: "Yen Press"},
		}

		source, err := yen.New(client)
		So(err, ShouldBeNil)

		firstCall := client.On("Do", mock.MatchedBy(matchFunc(id+"?next_ord=999"))).Maybe()
		secondCall := client.On("Do", mock.MatchedBy(matchFunc("?next_ord=4"))).Maybe()

		Convey("should return valid isbn", func() {
			firstCall.Return(openFile(t, "page_01.xml")).Once()
			secondCall.Return(openFile(t, "page_02.xml")).Once()

			books, err := source.GetISBNs(ctx, series)

			So(err, ShouldBeNil)
			So(books, ShouldResemble, expected)
		})
		Convey("should return error when", func() {
			Convey("making the initial call", func() {
				Convey("new request fails to make", func() {
					books, err := source.GetISBNs(ctx, types.Series{ID: string([]byte{0x7f})}) // this triggers failed parse

					So(err, ShouldBeError, `parse "https://yenpress.com/series/get_more/\x7f": net/url: invalid control character in URL`)
					So(books, ShouldBeNil)
				})
				Convey("do request fails", func() {
					firstCall.Return(nil, errors.New("do request failed")).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, "do request failed")
					So(books, ShouldBeEmpty)
				})
				Convey("html parse fails", func() {
					bodyMock.On("Read", mock.Anything).Return(0, errors.New("everybody body mock")).Once()
					firstCall.Return(bodyMock, nil).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, "everybody body mock")
					So(books, ShouldBeEmpty)
				})
			})
			Convey("making second call", func() {
				firstCall.Return(openFile(t, "page_01.xml")).Once()
				Convey("do request fails", func() {
					secondCall.Return(nil, errors.New("do request failed")).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, "do request failed")
					So(books, ShouldBeEmpty)
				})
				Convey("html parse fails", func() {
					bodyMock.
						On("Read", mock.Anything).
						Return(0, errors.New("everybody body mock")).
						Once()
					secondCall.Return(bodyMock, nil).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, "everybody body mock")
					So(books, ShouldBeEmpty)
				})
			})
		})
		Convey("should fail to sort when", func() {
			Convey("volume is not a number", func() {
				firstCall.Return(openFile(t, "page_01_bad_volume.xml")).Once()
				secondCall.Return(openFile(t, "page_02.xml")).Once()

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(books, ShouldNotResemble, expected)
			})
		})
		Convey("should blacklist isbn", func() {
			series.ISBNBlacklist = []string{expected[4].ISBN13, expected[11].ISBN13}
			firstCall.Return(openFile(t, "page_01.xml")).Once()
			secondCall.Return(openFile(t, "page_02.xml")).Once()

			books, err := source.GetISBNs(ctx, series)

			So(err, ShouldBeNil)
			So(books, ShouldResemble, append(append(expected[0:4], expected[5:11]...), expected[12:]...))
			So(books, ShouldNotResemble, expected)
		})
		Convey("should fail to parse book when", func() {
			Convey("the xml is missing an href", func() {
				firstCall.Return(openFile(t, "page_01_missing_href.xml")).Once()
				secondCall.Return(openFile(t, "page_02.xml")).Once()

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(books, ShouldResemble, expected[:18])
			})
			Convey("volume url is malformed", func() {
				expected := types.ISBNBooks{expected[18]}
				expected[0].Volume = ""
				expected[0].Title = " Vol. "
				firstCall.Return(openFile(t, "page_01_bad_volume_url.xml")).Once()

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(books, ShouldResemble, expected)
			})
		})
		Convey("should skip page 2 when data url", func() {
			Convey("is missing", func() {
				firstCall.Return(openFile(t, "page_01_missing_data_url.xml")).Once()

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(books, ShouldResemble, expected[4:])
			})
			Convey("contains invalid character", func() {
				firstCall.Return(openFile(t, "page_01_bad_data_url.xml")).Once()

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeError, `parse "/series/get_more/is-it-wrong-to-try-to-pick-up-girls-in-a-dungeon-light-novel?format%5B1%5D=Hardback&next_ord=66\x7f": net/url: invalid control character in URL`)
				So(books, ShouldBeNil)
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

	if !strings.HasSuffix(wd, filepath.Join("source", "yen")) {
		return os.Open(filepath.Join(wd, "source", "yen", "test_fixtures", fileName))
	}

	return os.Open(filepath.Join(wd, "test_fixtures", fileName))
}
