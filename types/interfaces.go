package types

import (
	"context"
)

type (
	SourceType  string
	LibraryType string
)

const (
	LibIBLibrary LibraryType = "libib"
	VizSource    SourceType  = "viz"
)

type ILibrary interface {
	GetBooksInCollection() ([]ISBNBook, error)
	SaveWanted(wanted []ISBNBook, withTitle bool) error
}

type ISource interface {
	GetISBNs(ctx context.Context, series Series) ([]ISBNBook, error)
}
