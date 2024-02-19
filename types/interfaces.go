package types

import (
	"context"

	"github.com/kevinanthony/collection-keep-updater/config"
)

type CollectionLibrary interface {
	GetBooksInCollection() ([]ISBNBook, error)
	SaveWanted(savePath string, wanted []ISBNBook, withTitle bool) error
}

type CollectionSource interface {
	GetISBNs(ctx context.Context, series config.Series) ([]ISBNBook, error)
}
