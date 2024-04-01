package wikipedia_test

import (
	"context"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	t.Parallel()

	Convey("New", t, func() {
		client := http.NewClientMock(t)
		getter := wikipedia.NewTableGetterMock(t)

		Convey("should return isource when http client is valid", func() {
			source, err := wikipedia.New(client, getter)

			So(source, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return error when", func() {
			Convey("http client is nil", func() {
				source, err := wikipedia.New(nil, getter)

				So(err, ShouldBeError, "http client is nil")
				So(source, ShouldBeNil)
			})
			Convey("table getter is nil", func() {
				source, err := wikipedia.New(client, nil)

				So(err, ShouldBeError, "wikipedia table getter is nil")
				So(source, ShouldBeNil)
			})
		})
	})
}

func TestWikiSource_GetISBNs(t *testing.T) {
	t.Parallel()

	Convey("GetISBNs", t, func() {
		ctx := context.Background()
		client := http.NewClientMock(t)
		getter := wikipedia.NewTableGetterMock(t)

		cfg := map[string]interface{}{
			"volume_header": "v_header",
			"title_header":  "t_header",
			"isbn_header":   "i_header",
			"tables":        []interface{}{81},
		}

		source, err := wikipedia.New(client, getter)

		series := types.Series{ID: "one piece", Name: "sanji is mid", SourceSettings: source.SourceSettingFromConfig(cfg)}

		tableCall := getter.On("GetTablesKeyValue", ctx, series.ID, "en", false, 1, 81).Maybe()

		t1 := []map[string]string{
			{"v_header": "1", "t_header": "vol 1", "i_header": "978-1-626923-48-5"},
		}
		t1_books := types.ISBNBooks{types.ISBNBook{Volume: "001", Title: "vol 1", ISBN13: "9781626923485", Source: "Wikipedia"}}
		t2 := []map[string]string{
			{"v_header": "2", "t_header": "vol 2", "i_header": "978-1-626924-31-4"},
			{"v_header": "3", "t_header": "vol 3", "i_header": "978-1-626924-85-7"},
			{"v_header": "4", "t_header": "vol 4", "i_header": "978-1-626925-46-5"},
		}
		t2_books := types.ISBNBooks{
			types.ISBNBook{Volume: "002", Title: "vol 2", ISBN13: "9781626924314", Source: "Wikipedia"},
			types.ISBNBook{Volume: "003", Title: "vol 3", ISBN13: "9781626924857", Source: "Wikipedia"},
			types.ISBNBook{Volume: "004", Title: "vol 4", ISBN13: "9781626925465", Source: "Wikipedia"},
		}

		So(err, ShouldBeNil)
		Convey("should return ISBNs", func() {
			Convey("when pages contains ISBN10 format", func() {
				t1 := []map[string]string{
					{"v_header": "1", "t_header": "vol 1", "i_header": "0123456789"},
					{"v_header": "2", "t_header": "vol 2", "i_header": "1234567890"},
				}
				t2 := []map[string]string{
					{"v_header": "3", "t_header": "vol 3", "i_header": "2345678901"},
					{"v_header": "4", "t_header": "vol 4", "i_header": "3456789012"},
				}

				expected := types.ISBNBooks{
					types.ISBNBook{Volume: "001", Title: "vol 1", ISBN10: "0123456789", Source: "Wikipedia"},
					types.ISBNBook{Volume: "002", Title: "vol 2", ISBN10: "1234567890", Source: "Wikipedia"},
					types.ISBNBook{Volume: "003", Title: "vol 3", ISBN10: "2345678901", Source: "Wikipedia"},
					types.ISBNBook{Volume: "004", Title: "vol 4", ISBN10: "3456789012", Source: "Wikipedia"},
				}

				tableCall.Once().Return([][]map[string]string{t1, t2}, nil)

				actual, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
			Convey("when pages contains ISBN13 format", func() {
				expected := append(t1_books, t2_books...)

				tableCall.Once().Return([][]map[string]string{t1, t2}, nil)

				actual, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
		})
		Convey("should exclude blacklisted ISBNs", func() {
			series.ISBNBlacklist = []string{"9781626924857"}

			expected := append(t1_books, t2_books[0], t2_books[2])

			tableCall.Once().Return([][]map[string]string{t1, t2}, nil)

			actual, err := source.GetISBNs(ctx, series)

			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})
		Convey("should return error when", func() {
			Convey("get settings fails", func() {
				series.SourceSettings = types.NewISourceSettingsMock(t)

				actual, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeError, "setting type not correct")
				So(actual, ShouldBeNil)
			})
			Convey("get tables key value fails", func() {
				tableCall.Once().Return(nil, errors.New("table parse failed"))

				actual, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeError, "table parse failed")
				So(actual, ShouldBeNil)
			})
		})
		Convey("should skip bad data when", func() {
			Convey("table setting volume header is nil", func() {
				cfg["volume_header"] = nil
				series.SourceSettings = source.SourceSettingFromConfig(cfg)
				t1_books[0].Volume = ""
				t2_books[0].Volume = ""
				t2_books[1].Volume = ""
				t2_books[2].Volume = ""

				tableCall.Once().Return([][]map[string]string{t1, t2}, nil)

				actual, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(actual, ShouldResemble, append(t1_books, t2_books...))
			})
			Convey("row does not have volume header", func() {
				t1_books[0].Volume = ""

				t1[0] = map[string]string{"junk": "1", "t_header": "vol 1", "i_header": "978-1-626923-48-5"}

				tableCall.Once().Return([][]map[string]string{t1}, nil)

				actual, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(actual, ShouldResemble, t1_books)
			})
			Convey("volume row is not a number", func() {
				t1_books[0].Volume = ""

				t1[0] = map[string]string{"v_header": "junk", "t_header": "vol 1", "i_header": "978-1-626923-48-5"}

				tableCall.Once().Return([][]map[string]string{t1}, nil)

				actual, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(actual, ShouldResemble, t1_books)
			})
			Convey("should composite title when", func() {
				Convey("table setting title header is nil", func() {
					cfg["title_header"] = nil
					series.SourceSettings = source.SourceSettingFromConfig(cfg)
					t1_books[0].Title = "sanji is mid Vol 001"
					t2_books[0].Title = "sanji is mid Vol 002"
					t2_books[1].Title = "sanji is mid Vol 003"
					t2_books[2].Title = "sanji is mid Vol 004"

					tableCall.Once().Return([][]map[string]string{t1, t2}, nil)

					actual, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeNil)
					So(actual, ShouldResemble, append(t1_books, t2_books...))
				})
				Convey("row does not have title header", func() {
					cfg["title_header"] = "nami-swan, wth"
					series.SourceSettings = source.SourceSettingFromConfig(cfg)
					t1_books[0].Title = "sanji is mid Vol 001"
					t2_books[0].Title = "sanji is mid Vol 002"
					t2_books[1].Title = "sanji is mid Vol 003"
					t2_books[2].Title = "sanji is mid Vol 004"

					tableCall.Once().Return([][]map[string]string{t1, t2}, nil)

					actual, err := source.GetISBNs(ctx, series)

					So(err, ShouldBeNil)
					So(actual, ShouldResemble, append(t1_books, t2_books...))
				})
			})
			Convey("table setting isbn header is nil", func() {
				cfg["isbn_header"] = nil
				series.SourceSettings = source.SourceSettingFromConfig(cfg)

				getter.On("GetTablesKeyValue", ctx, series.ID, "en", false, 1).
					Once().
					Return([][]map[string]string{t1, t2}, nil)

				actual, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(actual, ShouldBeEmpty)
			})
			Convey("row does not have isbn header", func() {
				cfg["isbn_header"] = "isbn"
				series.SourceSettings = source.SourceSettingFromConfig(cfg)

				tableCall.Once().Return([][]map[string]string{t1, t2}, nil)

				actual, err := source.GetISBNs(ctx, series)

				So(err, ShouldBeNil)
				So(actual, ShouldBeEmpty)
			})
		})
	})
}
