package main

import (
	"context"

	"github.com/kevinanthony/collection-keep-updater/collection/libib"
	"github.com/kevinanthony/collection-keep-updater/config"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
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

	libibSvc := libib.NewLibIB(cfg, httpClient)
	wikiDownloader := wikipedia.NewDownloader(httpClient)

	updateSvc := updater.New(libibSvc, wikiDownloader)

	if err := run(ctx, cfg, updateSvc); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, cfg config.App, updateSvc updater.Updater) error {
	allBooks, err := updateSvc.GetAllBooksForSeries(ctx, cfg.Series)
	if err != nil {
		return err
	}

	ownedBooks, err := updateSvc.GetLibraryBook(ctx)
	if err != nil {
		return err
	}

	wanted, err := updateSvc.Subtraction(ctx, ownedBooks, allBooks)
	if err != nil {
		return err
	}

	return updateSvc.SaveWanted(wanted)
}
