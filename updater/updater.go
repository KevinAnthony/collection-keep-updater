package updater

import (
	"context"
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/pkg/errors"
)

type Updater interface {
	GetAllAvailableBooks(ctx context.Context, series []types.Series) (types.ISBNBooks, error)
	UpdateLibrary(ctx context.Context, library types.ILibrary, availableBooks types.ISBNBooks) error
}

type updater struct {
	source map[types.SourceType]types.ISource
}

func New(source map[types.SourceType]types.ISource) Updater {
	return updater{
		source: source,
	}
}

func (u updater) UpdateLibrary(_ context.Context, library types.ILibrary, availableBooks types.ISBNBooks) error {
	booksInLibrary, err := library.GetBooksInCollection()
	if err != nil {
		return err
	}

	wanted, err := booksInLibrary.Diff(availableBooks)
	if err != nil {
		return err
	}

	err = u.SaveWanted(library, wanted)
	if err != nil {
		return err
	}

	return nil
}

func (u updater) GetAllAvailableBooks(ctx context.Context, series []types.Series) (types.ISBNBooks, error) {
	allBooks := types.NewISBNBooks(0)

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

func (u updater) SaveWanted(library types.ILibrary, wanted types.ISBNBooks) error {
	if len(wanted) == 0 {
		fmt.Println("No New Wanted books")

		return nil
	}

	return library.SaveWanted(wanted, false)
}
