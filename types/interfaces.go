package types

import (
	"context"
)

type (
	SourceType  string
	LibraryType string
)

type ILibrary interface {
	GetBooksInCollection() (ISBNBooks, error)
	SaveWanted(wanted ISBNBooks, withTitle bool) error
}

type ISource interface {
	GetISBNs(ctx context.Context, series Series) (ISBNBooks, error)
}
