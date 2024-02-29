package types

import (
	"context"

	"github.com/spf13/cobra"
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
	ISourceSetting
}

type ISourceSetting interface {
	SetFlags(cmd *cobra.Command)
	SourceSettingFromConfig(data map[string]interface{}) ISourceSettings
	SourceSettingFromFlags(cmd *cobra.Command, original ISourceSettings) (ISourceSettings, error)
	GetIDFromURL(url string) (string, error)
}
