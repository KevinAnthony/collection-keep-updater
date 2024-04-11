package viz_test

import (
	"bytes"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVizSettings_Print(t *testing.T) {
	t.Parallel()

	Convey("Print", t, func() {
		writer := bytes.NewBufferString("")
		cmdMock := types.NewICommandMock(t)
		settings := getSettings(t)

		cmdMock.On("OutOrStdout").Once().Return(writer)

		Convey("should print empty message to buffer", func() {
			err := settings.Print(cmdMock)

			So(err, ShouldBeNil)

			So(writer.String(), ShouldEqual, "┌─────────────────┬───────┐\n│ MAXIMUM BACKLOG │ DELAY │\n├─────────────────┼───────┤\n│               3 │  50ms │\n└─────────────────┴───────┘\n")
		})
	})
}

func getSettings(t *testing.T) types.ISourceSettings {
	t.Helper()

	source := getSource(t)

	data := map[string]interface{}{
		"maximum_backlog": 3,
		"delay_between":   "50ms",
	}

	return source.SourceSettingFromConfig(data)
}

func getSource(t *testing.T) types.ISource {
	t.Helper()

	cmd := types.NewICommandMock(t)
	ctx := ctxu.NewContextMock(t)
	client := http.NewClientMock(t)

	cmd.On("Context").Return(ctx)
	ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(client)

	return viz.New(cmd)
}
