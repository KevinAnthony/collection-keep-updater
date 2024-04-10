package wikipedia_test

import (
	"bytes"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"

	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWikiSettings_Print(t *testing.T) {
	t.Parallel()

	table := `┌────────────────────┬───────────────────┬───────────────────┬────────────────┐
│ VOLUME COLUMN NAME │ TITLE COLUMN NAME │ ISBN COLUMN TITLE │ TABLES ON PAGE │
├────────────────────┼───────────────────┼───────────────────┼────────────────┤
│ vheader            │ theader           │ iheader           │ 1, 8, 12       │
└────────────────────┴───────────────────┴───────────────────┴────────────────┘
`

	Convey("Print", t, func() {
		writer := bytes.NewBufferString("")
		cmdMock := types.NewICommandMock(t)
		settings := getSettings(t)

		cmdMock.On("OutOrStdout").Once().Return(writer)

		Convey("should print empty message to buffer", func() {
			settings.Print(cmdMock)

			So(writer.String(), ShouldEqual, table)
		})
	})
}

func getSettings(t *testing.T) types.ISourceSettings {
	t.Helper()

	source := getSource(t)

	data := map[string]interface{}{
		"volume_header": "vheader",
		"title_header":  "theader",
		"isbn_header":   "iheader",
		"tables":        []interface{}{1, 8, 12},
	}

	return source.SourceSettingFromConfig(data)
}

func getSource(t *testing.T) types.ISource {
	t.Helper()

	cmd := types.NewICommandMock(t)
	ctx := ctxu.NewContextMock(t)
	client := http.NewClientMock(t)
	getter := wikipedia.NewTableGetterMock(t)

	cmd.On("Context").Return(ctx)
	ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(client)
	ctx.On("Value", ctxu.ContextKey("wiki_getter_ctx_key")).Return(getter)

	return wikipedia.New(cmd)
}
