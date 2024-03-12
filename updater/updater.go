package updater

import (
	"context"

	"github.com/kevinanthony/collection-keep-updater/ctxu"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/pkg/errors"
)

//go:generate mockery --name=IUpdater --structname=IUpdaterMock --filename=updater_mock.go --inpackage
type IUpdater interface {
	GetAllAvailableBooks(ctx types.ICommand, series []types.Series) (types.ISBNBooks, error)
	UpdateLibrary(ctx context.Context, library types.ILibrary, availableBooks types.ISBNBooks) (types.ISBNBooks, error)
}

type updater struct {
}

func NewUpdater() IUpdater {
	return updater{}
}

func (u updater) UpdateLibrary(_ context.Context, library types.ILibrary, availableBooks types.ISBNBooks) (types.ISBNBooks, error) {
	booksInLibrary, err := library.GetBooksInCollection()
	if err != nil {
		return nil, err
	}

	wanted, err := booksInLibrary.Diff(availableBooks)
	if err != nil {
		return nil, err
	}

	return wanted, nil
}

func (u updater) GetAllAvailableBooks(cmd types.ICommand, series []types.Series) (types.ISBNBooks, error) {
	allBooks := types.NewISBNBooks(0)

	for _, s := range series {
		if len(s.ID) == 0 {
			continue
		}

		source, err := ctxu.GetSource(cmd, s.Source)
		if err != nil {
			return nil, errors.Wrapf(err, "%s is unknown", s.Source)
		}

		books, err := source.GetISBNs(cmd.Context(), s)
		if err != nil {
			return nil, errors.Wrapf(err, s.Name)
		}

		allBooks = append(allBooks, books...)
	}

	return allBooks, nil
}
