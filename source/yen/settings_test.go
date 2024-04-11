package yen_test

import (
	"bytes"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/source/yen"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
)

func TestYenSettings_Print(t *testing.T) {
	t.Parallel()

	Convey("Print", t, func() {
		writer := bytes.NewBufferString("")
		cmdMock := types.NewICommandMock(t)
		settings := getSettings(t)

		cmdMock.On("OutOrStdout").Once().Return(writer)

		Convey("should print empty message to buffer", func() {
			err := settings.Print(cmdMock)

			So(err, ShouldBeNil)

			So(writer.String(), ShouldEqual, "┌─────────────────┐\n│ NO YEN SETTINGS │\n├─────────────────┤\n└─────────────────┘\n")
		})
	})
}

func getSettings(t *testing.T) types.ISourceSettings {
	t.Helper()

	cmd := types.NewICommandMock(t)
	ctx := ctxu.NewContextMock(t)
	client := http.NewClientMock(t)

	cmd.On("Context").Return(ctx)
	ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(client)

	source := yen.New(cmd)

	return source.SourceSettingFromConfig(nil)
}
