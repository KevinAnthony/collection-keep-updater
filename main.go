package main

import (
	"context"

	"github.com/kevinanthony/collection-keep-updater/config"
	"github.com/kevinanthony/collection-keep-updater/library/libib"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/updater"
	"github.com/kevinanthony/gorps/v2/encoder"
	"github.com/kevinanthony/gorps/v2/http"
)

func main() {
	ctx := context.Background()

	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	httpClient := http.NewClient(http.NewNativeClient(), encoder.NewFactory())

	libraries := map[types.LibraryType]types.ILibrary{}
	for _, setting := range cfg.Libraries {
		switch setting.Name {
		case types.LibIBLibrary:
			libraries[types.LibIBLibrary] = libib.New(setting, httpClient)
		}
	}
	sources := map[types.SourceType]types.ISource{
		types.WikipediaSource: wikipedia.New(httpClient),
		types.VizSource:       viz.New(httpClient),
	}

	updateSvc := updater.New(sources)

	if err := run(ctx, cfg, libraries, updateSvc); err != nil {
		panic(err)
	}
}

func run(
	ctx context.Context,
	cfg config.App,
	libraries map[types.LibraryType]types.ILibrary,
	updateSvc updater.Updater,
) error {
	availableBooks, err := updateSvc.GetAllAvailableBooks(ctx, cfg.Series)
	if err != nil {
		return err
	}

	for _, library := range libraries {
		err := updateSvc.UpdateLibrary(ctx, library, availableBooks)
		if err != nil {
			return err
		}
	}

	return nil
}
