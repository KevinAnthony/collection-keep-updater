package types_test

import (
	"bytes"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSeries_Print(t *testing.T) {
	t.Parallel()

	configTable := "┌──────────┬───────────┬─────────────┬─────────┐\n│ KEY      │ NAME      │ SOURCE      │ ID      │\n├──────────┼───────────┼─────────────┼─────────┤\n│ test key │ test name │ test source │ test ID │\n└──────────┴───────────┴─────────────┴─────────┘\n"
	isbnTable := "┌─────────────────┐\n│ ISBN BLACKLIST  │\n├─────────────────┤\n│ blacklist eps 1 │\n│ blacklist eps 2 │\n└─────────────────┘\n"
	Convey("Print", t, func() {
		s := types.Series{
			Name:   "test name",
			ID:     "test ID",
			Source: "test source",
			Key:    "test key",
		}
		writer := bytes.NewBufferString("")
		cmdMock := types.NewICommandMock(t)
		sourceSettingMock := types.NewISourceSettingsMock(t)

		outCall := cmdMock.On("OutOrStdout").Maybe()
		sourcePrintCall := sourceSettingMock.On("Print", cmdMock).Maybe()

		Convey("should print the list of series", func() {
			Convey("if source settings is nil", func() {
				outCall.Once().Return(writer)

				err := s.Print(cmdMock)

				So(err, ShouldBeNil)
				So(writer.String(), ShouldEqual, configTable)
			})
			Convey("if blacklist is not nil", func() {
				outCall.Twice().Return(writer)
				s.ISBNBlacklist = []string{"blacklist eps 1", "blacklist eps 2"}

				err := s.Print(cmdMock)

				So(err, ShouldBeNil)
				So(writer.String(), ShouldEqual, configTable+isbnTable)
			})
			Convey("and source settings if source settings is not nil", func() {
				s.SourceSettings = sourceSettingMock
				outCall.Once().Return(writer)
				Convey("and source print returns nil", func() {
					sourcePrintCall.Once().Return(nil)

					err := s.Print(cmdMock)

					So(err, ShouldBeNil)
					So(writer.String(), ShouldEqual, configTable)
				})
				Convey("and source print returns an error", func() {
					sourcePrintCall.Once().Return(errors.New("test error"))

					err := s.Print(cmdMock)

					So(err, ShouldBeError, "test error")
					So(writer.String(), ShouldEqual, configTable)
				})
			})

			cmdMock.AssertExpectations(t)
			sourceSettingMock.AssertExpectations(t)
		})
	})
}

func TestSeries_String(t *testing.T) {
	t.Parallel()

	Convey("String", t, func() {
		Convey("should return string", func() {
			s := types.Series{
				Name:   "test name",
				ID:     "test ID",
				Source: "test source",
				Key:    "test key",
			}

			So(s.String(), ShouldEqual, "test name (test source)")
		})
	})
}

func TestGetSetting(t *testing.T) {
	t.Parallel()

	Convey("GetSetting", t, func() {
		sourceSettingMock := types.NewISourceSettingsMock(t)

		s := types.Series{
			Name:   "test name",
			ID:     "test ID",
			Source: "test source",
			Key:    "test key",
		}

		Convey("if source settings is set and correct", func() {
			s.SourceSettings = sourceSettingMock

			m, err := types.GetSetting[*types.ISourceSettingsMock](s)

			So(m, ShouldEqual, sourceSettingMock)
			So(err, ShouldBeNil)
		})
		Convey("should return nil, nil if source settings is nil", func() {
			m, err := types.GetSetting[*types.ISourceSettingsMock](s)

			So(m, ShouldBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return error if if source setting is not correct", func() {
			s.SourceSettings = sourceSettingMock

			m, err := types.GetSetting[testSourceSetting](s)

			So(err, ShouldBeError, "setting type not correct")
			So(m, ShouldResemble, testSourceSetting{})
		})
	})
}

type testSourceSetting struct{}

func (t testSourceSetting) Print(_ types.ICommand) error { return nil }
