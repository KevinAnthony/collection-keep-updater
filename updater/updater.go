package updater

import (
	"context"
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/config"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"
)

type Updater interface {
	GetAllBooksForSeries(ctx context.Context, series []config.Series) ([]types.ISBNBook, error)
	GetLibraryBook(ctx context.Context) ([]types.ISBNBook, error)
	Subtraction(ctx context.Context, library, all []types.ISBNBook) ([]types.ISBNBook, error)
	SaveWanted(wanted []types.ISBNBook) error
}

type updater struct {
	source   types.CollectionSource
	library  types.CollectionLibrary
	savePath string
}

func New(library types.CollectionLibrary, source types.CollectionSource) Updater {
	return updater{
		library:  library,
		source:   source,
		savePath: "wanted.csv",
	}
}

func (u updater) Subtraction(_ context.Context, library, allPublished []types.ISBNBook) ([]types.ISBNBook, error) {
	var wanted []types.ISBNBook
	for _, book := range allPublished {
		if utils.Contains(library, book, types.ISBNBookCmp) {
			continue
		}

		wanted = append(wanted, book)
	}

	return wanted, nil
}

func (u updater) GetLibraryBook(_ context.Context) ([]types.ISBNBook, error) {
	return u.library.GetBooksInCollection()
}

func (u updater) GetAllBooksForSeries(ctx context.Context, series []config.Series) ([]types.ISBNBook, error) {
	var allBooks []types.ISBNBook

	for _, s := range series {
		if len(s.ID) == 0 {
			continue
		}

		books, err := u.source.GetISBNs(ctx, s)
		if err != nil {
			return nil, err
		}

		allBooks = append(allBooks, books...)
	}

	// printList(allBooks)

	return allBooks, nil
}

func (u updater) SaveWanted(wanted []types.ISBNBook) error {
	if len(wanted) == 0 {
		fmt.Println("No New Wanted books")

		return nil
	}

	return u.library.SaveWanted(u.savePath, wanted, false)
}
