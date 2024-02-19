package updater

import (
	"context"
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/pkg/errors"
)

type Updater interface {
	GetAllAvailableBooks(ctx context.Context, series []types.Series) ([]types.ISBNBook, error)
	UpdateLibrary(ctx context.Context, library types.ILibrary, availableBooks []types.ISBNBook) error
}

type updater struct {
	source map[types.SourceType]types.ISource
}

func New(source map[types.SourceType]types.ISource) Updater {
	return updater{
		source: source,
	}
}

func (u updater) UpdateLibrary(_ context.Context, library types.ILibrary, availableBooks []types.ISBNBook) error {
	booksInLibrary, err := library.GetBooksInCollection()
	if err != nil {
		return err
	}

	wanted, err := subtraction(booksInLibrary, availableBooks)
	if err != nil {
		return err
	}

	err = u.SaveWanted(library, wanted)
	if err != nil {
		return err
	}

	return nil
}

func (u updater) GetAllAvailableBooks(ctx context.Context, series []types.Series) ([]types.ISBNBook, error) {
	var allBooks []types.ISBNBook

	for _, s := range series {
		if len(s.ID) == 0 {
			continue
		}

		source, found := u.source[s.Source]
		if !found {
			return nil, fmt.Errorf("source: %s is unknown", s.Source)
		}
		books, err := source.GetISBNs(ctx, s)
		if err != nil {
			return nil, errors.Wrap(err, s.Name)
		}

		allBooks = append(allBooks, books...)
	}

	return allBooks, nil
}

func (u updater) SaveWanted(library types.ILibrary, wanted []types.ISBNBook) error {
	if len(wanted) == 0 {
		fmt.Println("No New Wanted books")

		return nil
	}

	return library.SaveWanted(wanted, false)
}

func subtraction(minuend, subtrahend []types.ISBNBook) ([]types.ISBNBook, error) {
	var diff []types.ISBNBook
	for _, book := range subtrahend {
		if utils.Contains(minuend, book, types.ISBNBookCmp) {
			continue
		}

		diff = append(diff, book)
	}

	return diff, nil
}
