package kodansha_test

import (
	"bytes"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/source/kodansha"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/kevinanthony/collection-keep-updater/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKondashaSettings_Print(t *testing.T) {
	t.Parallel()

	Convey("Print", t, func() {
		writer := bytes.NewBufferString("")
		cmdMock := types.NewICommandMock(t)
		settings := getSettings(t)

		cmdMock.On("OutOrStdout").Once().Return(writer)

		Convey("should print empty message to buffer", func() {
			err := settings.Print(cmdMock)

			So(err, ShouldBeNil)

			So(writer.String(), ShouldEqual, "┌──────────────────────┐\n│ NO KONDASHA SETTINGS │\n├──────────────────────┤\n└──────────────────────┘\n")
		})
	})
}

func getSettings(t *testing.T) types.ISourceSettings {
	t.Helper()

	client := http.NewClientMock(t)

	source, err := kodansha.New(client)
	So(err, ShouldBeNil)

	return source.SourceSettingFromConfig(nil)
}
