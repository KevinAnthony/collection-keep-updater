package libib

import (
	"context"
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/gocarina/gocsv"
)

const (
	exportURL   = "https://www.libib.com/settings/export-library/submit"
	validateURL = "https://www.libib.com/csvimport/validate-file/submit"
	processURL  = "https://www.libib.com/csvimport/process-import"
)

type libIB struct {
	cfg    types.LibrarySettings
	client http.Client
	helper ILibibHelper
}

func New(cmd types.ICommand, cfg types.LibrarySettings) types.ILibrary {
	return libIB{
		client: ctxu.GetHttpClient(cmd),
		cfg:    cfg,
		helper: NewLibinHelper(cmd, cfg),
	}
}

func (l libIB) GetBooksInCollection(ctx context.Context) (types.ISBNBooks, error) {
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

func (l libIB) SaveWanted(cmd types.ICommand, wanted types.ISBNBooks) error {
	results, err := l.helper.SubmitForm(cmd, wanted)
	if err != nil {
		return err
	}

	isbn, err := l.helper.GetISBNFromSuccessResult(cmd, results)
	if err != nil {
		return err
	}

	if err := l.helper.ValidateResults(cmd, wanted, isbn); err != nil {
		return err
	}

	cols, err := l.helper.GetQueryParamsFromSuccessResults(cmd, results)
	if err != nil {
		return err
	}

	return l.helper.SaveResults(cmd, cols)
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

func (l libIB) getCSV(ctx context.Context, libraryID string) ([]libibCSVEntries, error) {
	req, _ := http.
		NewRequest(l.client).
		Post().
		URL(exportURL).
		Body(fmt.Sprintf("settings-library-export-id=%s", libraryID)).
		Header("Cookie", l.cfg.APIKey).
		Header("Content-Type", "application/x-www-form-urlencoded").
		CreateRequest(ctx)

	body, err := l.client.Do(req)
	if err != nil {
		// TODO: log error
		return nil, err
	}

	var library []libibCSVEntries

	if err := gocsv.Unmarshal(body, &library); err != nil { // Load clients from file
		return nil, err
	}

	return library, nil
}
