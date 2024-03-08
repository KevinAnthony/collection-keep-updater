package libib

import (
	"context"
	"fmt"
	"os"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/gocarina/gocsv"
)

const (
	exportURL   = "https://www.libib.com/settings/export-library/submit"
	outFileName = "wanted.csv"
)

type libIB struct {
	cfg    types.LibrarySettings
	client http.Client
}

func New(cfg types.LibrarySettings, c http.Client) types.ILibrary {
	if c == nil {
		panic("http client is nil")
	}

	return libIB{
		client: c,
		cfg:    cfg,
	}
}

func (l libIB) GetBooksInCollection() (types.ISBNBooks, error) {
	ctx := context.Background()
	var libibEntries []libibCSVEntries

	for _, library := range append(l.cfg.OtherColIDs, l.cfg.WantedColID) {
		entries, err := l.getCSV(ctx, library)
		if err != nil {
			return nil, err
		}
		libibEntries = append(libibEntries, entries...)
	}

	return l.createISBNBook(libibEntries), nil
}

func (l libIB) SaveWanted(wanted types.ISBNBooks) error {
	outFile, err := os.Create(outFileName)
	if err != nil {
		return err
	}

	return gocsv.MarshalFile(l.createCSVEntries(wanted), outFile)
}

func (l libIB) createISBNBook(entries []libibCSVEntries) types.ISBNBooks {
	books := types.NewISBNBooks(len(entries))
	for _, entry := range entries {
		books = append(books, types.ISBNBook{
			ISBN10: entry.ISBN,
			ISBN13: entry.ISBN13,
			Title:  entry.Title,
		})
	}

	return books
}

func (l libIB) createCSVEntries(books types.ISBNBooks) []libibCSVEntries {
	entries := make([]libibCSVEntries, 0, len(books))
	for _, book := range books {
		entry := libibCSVEntries{
			ISBN13: book.ISBN13,
			ISBN:   book.ISBN10,
		}

		entries = append(entries, entry)
	}

	return entries
}

func (l libIB) getCSV(ctx context.Context, libraryID string) ([]libibCSVEntries, error) {
	req, err := http.
		NewRequest(l.client).
		Post().
		URL(exportURL).
		Body(fmt.Sprintf("settings-library-export-id=%s", libraryID)).
		Header("Cookie", l.cfg.APIKey).
		Header("Content-Type", "application/x-www-form-urlencoded").
		CreateRequest(ctx)
	if err != nil {
		return nil, err
	}

	body, err := l.client.Do(req)
	if err != nil {
		// TODO: log error
		return nil, err
	}

	var library []libibCSVEntries

	if err := gocsv.UnmarshalBytes(body, &library); err != nil { // Load clients from file
		return nil, err
	}

	return library, nil
}
