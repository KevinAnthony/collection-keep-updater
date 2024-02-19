package types

import (
	"context"
)

type SourceType string

type CollectionLibrary interface {
	GetBooksInCollection() ([]ISBNBook, error)
	SaveWanted(savePath string, wanted []ISBNBook, withTitle bool) error
}

type CollectionSource interface {
	GetISBNs(ctx context.Context, series Series) ([]ISBNBook, error)
}
