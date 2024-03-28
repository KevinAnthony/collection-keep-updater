package viz_test

import (
	native "net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/stretchr/testify/mock"

	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	t.Parallel()

	Convey("New", t, func() {
		client := http.NewClientMock(t)
		Convey("should return isource when http client is valid", func() {
			source, err := viz.New(client)

			So(source, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return error when http client is nil", func() {
			source, err := viz.New(nil)

			So(err, ShouldBeError, "http client is nil")
			So(source, ShouldBeNil)
		})
	})
}

func TestViz_GetISBNs(t *testing.T) {
	t.Parallel()

	Convey("GetISBNs", t, func() {
		id := "chainsaw-man"
		client := http.NewClientMock(t)
		ctx := ctxu.NewContextMock(t)
		bodyMock := http.NewBodyMock(t)

		series := types.Series{ID: id, Name: "Chainsaw The Man"}
		expected := types.ISBNBooks{
			{
				ISBN13: "9781974709939",
				Volume: "1",
				Title:  "Chainsaw The Man: #1",
				Source: "Viz",
			},
			{
				ISBN13: "9781974709946",
				Volume: "2",
				Title:  "Chainsaw The Man: #2",
				Source: "Viz",
			},
			{
				ISBN13: "9781974709953",
				Volume: "3",
				Title:  "Chainsaw The Man: #3",
				Source: "Viz",
			},
			{
				ISBN13: "9781974717279",
				Volume: "4",
				Title:  "Chainsaw The Man: #4",
				Source: "Viz",
			},
		}
		source, err := viz.New(client)
		So(err, ShouldBeNil)

		seriesCall := client.On("Do", mock.MatchedBy(matchFunc(id+"/all"))).Maybe()
		page1Call := client.On("Do", mock.MatchedBy(matchFunc("6419"))).Maybe()
		page2Call := client.On("Do", mock.MatchedBy(matchFunc("6495"))).Maybe()
		page3Call := client.On("Do", mock.MatchedBy(matchFunc("6567"))).Maybe()
		page4Call := client.On("Do", mock.MatchedBy(matchFunc("6627"))).Maybe()

		Convey("should return valid isbn", func() {
			Convey("and all the pages when max backlog is nil", func() {
				seriesCall.Return(openFile(t, "series.html")).Once()
				page1Call.Return(openFile(t, "page_1.html")).Once()
				page2Call.Return(openFile(t, "page_2.html")).Once()
				page3Call.Return(openFile(t, "page_3.html")).Once()
				page4Call.Return(openFile(t, "page_4.html")).Once()

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(books, ShouldResemble, expected)
			})
			Convey("and only last n pages when max backlog is not nil", func() {
				series.SourceSettings = source.SourceSettingFromConfig(map[string]interface{}{"maximum_backlog": 3})
				expected = expected[1:]
				seriesCall.Return(openFile(t, "series.html")).Once()
				page2Call.Return(openFile(t, "page_2.html")).Once()
				page3Call.Return(openFile(t, "page_3.html")).Once()
				page4Call.Return(openFile(t, "page_4.html")).Once()

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(books, ShouldResemble, expected)
			})
			Convey("and should delay when delay is set", func() {
				series.SourceSettings = source.SourceSettingFromConfig(
					map[string]interface{}{
						"maximum_backlog": 3,
						"delay_between":   "50ms",
					},
				)
				expected = expected[1:]
				seriesCall.Return(openFile(t, "series.html")).Once()
				page2Call.Return(openFile(t, "page_2.html")).Once()
				page3Call.Return(openFile(t, "page_3.html")).Once()
				page4Call.Return(openFile(t, "page_4.html")).Once()

				start := time.Now()
				books, err := source.GetISBNs(ctx, series)
				end := time.Since(start)

				So(err, ShouldBeNil)
				So(books, ShouldResemble, expected)
				So(end, ShouldBeBetween, 150*time.Millisecond, 175*time.Millisecond)
			})
			Convey("should exclude books that have blacklists", func() {
				series.ISBNBlacklist = []string{expected[2].ISBN13}
				expected = types.ISBNBooks{expected[0], expected[1], expected[3]}

				seriesCall.Return(openFile(t, "series.html")).Once()
				page1Call.Return(openFile(t, "page_1.html")).Once()
				page2Call.Return(openFile(t, "page_2.html")).Once()
				page3Call.Return(openFile(t, "page_3.html")).Once()
				page4Call.Return(openFile(t, "page_4.html")).Once()

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(books, ShouldResemble, expected)
			})
		})
		Convey("should return error when", func() {
			Convey("source settings is not a viz setting", func() {
				series.SourceSettings = types.NewISourceSettingsMock(t)

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeError, `setting type not correct`)
				So(books, ShouldBeNil)
			})
			Convey("for series", func() {
				Convey("create request fails", func() {
					books, err := source.GetISBNs(ctx, types.Series{ID: string([]byte{0x7f})}) // this triggers failed parse

					So(err, ShouldBeError, `parse "https://www.viz.com/read/manga/\x7f/all": net/url: invalid control character in URL`)
					So(books, ShouldBeNil)
				})
				Convey("do request fails", func() {
					seriesCall.Return(nil, errors.New("do request error")).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, `do request error`)
					So(books, ShouldBeNil)
				})
				Convey("parse fails", func() {
					bodyMock.On("Read", mock.Anything).Return(0, errors.New("everybody body mock")).Once()
					seriesCall.Return(bodyMock, nil).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, "everybody body mock")
					So(books, ShouldBeNil)
				})
			})
			Convey("for page", func() {
				Convey("create request fails", func() {
					seriesCall.Return(openFile(t, "series_bad_url.html")).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, "/read/manga/chainsaw-man-volume-1/product/6419\u007F: parse \"https://www.viz.com//read/manga/chainsaw-man-volume-1/product/6419\\x7f\": net/url: invalid control character in URL")
					So(books, ShouldBeNil)
				})
				Convey("do request fails", func() {
					seriesCall.Return(openFile(t, "series.html")).Once()
					page1Call.Return(nil, errors.New("do request error")).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, "/read/manga/chainsaw-man-volume-1/product/6419: do request error")
					So(books, ShouldBeNil)
				})
				Convey("parse fails", func() {
					bodyMock.On("Read", mock.Anything).Return(0, errors.New("everybody body mock")).Once()
					seriesCall.Return(openFile(t, "series.html")).Once()
					page1Call.Return(bodyMock, nil).Once()

					books, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeError, "/read/manga/chainsaw-man-volume-1/product/6419: everybody body mock")
					So(books, ShouldBeNil)
				})
			})
		})
		Convey("should partially fail when", func() {
			Convey("strong node is missing it's children", func() {
				seriesCall.Return(openFile(t, "series.html")).Once()
				page1Call.Return(openFile(t, "page_1_empty_strong.html")).Once()
				page2Call.Return(openFile(t, "page_2.html")).Once()
				page3Call.Return(openFile(t, "page_3.html")).Once()
				page4Call.Return(openFile(t, "page_4.html")).Once()

				expected = expected[1:]

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(books, ShouldResemble, expected)
			})
		})
		Convey("parse volume from path fails because", func() {
			Convey("url has less the 4 slugs", func() {
				expected[0].Title = "Chainsaw The Man: #"
				expected[0].Volume = ""
				seriesCall.Return(openFile(t, "series_invalid_url.html")).Once()
				page1Call.Return(openFile(t, "page_1.html")).Once()
				page2Call.Return(openFile(t, "page_2.html")).Once()
				page3Call.Return(openFile(t, "page_3.html")).Once()
				page4Call.Return(openFile(t, "page_4.html")).Once()

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(books, ShouldResemble, expected)
			})
			Convey("volume numbers is not a number", func() {
				expected[0].Title = "Chainsaw The Man: #"
				expected[0].Volume = ""
				seriesCall.Return(openFile(t, "series_invalid_volume.html")).Once()
				page1Call.Return(openFile(t, "page_1.html")).Once()
				page2Call.Return(openFile(t, "page_2.html")).Once()
				page3Call.Return(openFile(t, "page_3.html")).Once()
				page4Call.Return(openFile(t, "page_4.html")).Once()

				books, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(books, ShouldResemble, expected)
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

	if !strings.HasSuffix(wd, filepath.Join("source", "viz")) {
		return os.Open(filepath.Join(wd, "source", "viz", "test_fixtures", fileName))
	}

	return os.Open(filepath.Join(wd, "test_fixtures", fileName))
}
