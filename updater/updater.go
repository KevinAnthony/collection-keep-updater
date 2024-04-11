package updater

import (
	"context"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/pkg/errors"
)

const updaterCtxKey ctxu.ContextKey = "updater_ctx_key"

//go:generate mockery --name=IUpdater --structname=IUpdaterMock --filename=updater_mock.go --inpackage
type IUpdater interface {
	GetAllAvailableBooks(cmd types.ICommand, series []types.Series) (types.ISBNBooks, error)
	UpdateLibrary(ctx context.Context, library types.ILibrary, availableBooks types.ISBNBooks) (types.ISBNBooks, error)
}

type updater struct{}

func NewUpdater() IUpdater {
	return updater{}
}

func GetUpdater(cmd types.ICommand) IUpdater {
	ctx := cmd.Context()

	value := ctx.Value(updaterCtxKey)
	if client, ok := value.(IUpdater); ok {
		return client
	}

	u := NewUpdater()
	ctx = context.WithValue(ctx, updaterCtxKey, u)
	cmd.SetContext(ctx)

	return u
}

func (u updater) UpdateLibrary(ctx context.Context, library types.ILibrary, availableBooks types.ISBNBooks) (types.ISBNBooks, error) {
	booksInLibrary, err := library.GetBooksInCollection(ctx)
	if err != nil {
		return nil, err
	}

	return availableBooks.Diff(booksInLibrary), nil
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
