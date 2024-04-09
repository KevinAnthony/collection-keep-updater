package cmd_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/cmd"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetRootCmd(t *testing.T) {
	t.Parallel()
	Convey("GetRootCmd", t, func() {
		Convey("should return non-nil command", func() {
			So(cmd.GetRootCmd(), ShouldNotBeNil)
		})
	})
}

//func TestPreREunE(t *testing.T) {
//	t.Parallel()
//
//	Convey("PreREunE", t, func() {
//		ctx := ctxu.NewContextMock(t)
//		sourceMock := types.NewISouceMock(t)
//		cfgMock := types.NewIConfigMock(t)
//		//source := types.NewISouceMock(t)
//
//		command := &cobra.Command{}
//		command.SetContext(ctx)
//		settingMap := map[string]any{
//			"delay_between":   "100ms",
//			"maximum_backlog": 2,
//		}
//		//settingMap := map[string]any{
//		//	"delay_between":   "100ms",
//		//	"maximum_backlog": 2,
//		//}
//		seriesBlob := []any{
//			map[string]any{
//				"id":              "one-piece",
//				"key":             "one-piece",
//				"name":            "One Piece",
//				"source":          "viz",
//				"source_settings": settingMap,
//				"isbn_blacklist":  []any{"one", "two", "five"},
//			},
//		}
//		libraryBlob := []any{
//			map[string]any{
//				"api_key":              "secret",
//				"other_collection_ids": []any{"id1", "id2", "id3"},
//				"wanted_collection_id": "id0",
//				"type":                 "libib",
//			},
//		}
//		vizSrc, err := viz.New(http.NewClientMock(t))
//		So(err, ShouldBeNil)
//
//		readCall := mockLoadCfg(cfgMock)
//
//		sourceSetting := vizSrc.SourceSettingFromConfig(settingMap)
//
//		sources := map[types.SourceType]types.ISource{types.VizSource: sourceMock}
//
//		getConfigCall := ctx.On("Value", ctxu.ContextKey("i_config_ctx_key")).Return(cfgMock).Maybe()
//		getSourceCall := ctx.On("Value", ctxu.ContextKey("sources_ctx_key")).Return(sources).Maybe()
//		sourceSettingsCall := sourceMock.On("SourceSettingFromConfig", settingMap).Maybe()
//		getSeriesCall := cfgMock.On("Get", "series").Maybe()
//		getLibCall := cfgMock.On("Get", "libraries").Maybe()
//		ctxCall := ctx.On("Value", ctxu.ContextKey("sources_ctx_key")).Maybe()
//
//		Convey("should return no errors", func() {
//			getConfigCall.Once()
//			getSourceCall.Once()
//			readCall.Return(nil)
//			getSeriesCall.Return(seriesBlob).Once()
//			ctxCall.Return(sources).Once()
//			sourceSettingsCall.Return(sourceSetting).Once()
//
//			getLibCall.Return(libraryBlob).Once()
//
//			err := cmd.PreREunE(command, nil)
//
//			So(err, ShouldBeNil)
//		})
//	})
//}
