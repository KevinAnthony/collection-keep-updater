package types

import (
	"context"

	"github.com/spf13/cobra"
)

type (
	SourceType  string
	LibraryType string
)

//go:generate mockery --name=ILibrary --structname=ILibraryMock --filename=library_mock.go --inpackage
type ILibrary interface {
	GetBooksInCollection() (ISBNBooks, error)
	SaveWanted(wanted ISBNBooks) error
}

//go:generate mockery --name=ISource --structname=ISouceMock --filename=source_mock.go --inpackage
type ISource interface {
	GetISBNs(ctx context.Context, series Series) (ISBNBooks, error)
	ISourceConfig
}

//go:generate mockery --name=ISourceConfig --structname=ISourceConfigMock --filename=source_config_mock.go --inpackage
type ISourceConfig interface {
	SourceSettingFromConfig(data map[string]interface{}) ISourceSettings
	SourceSettingFromFlags(cmd *cobra.Command, original ISourceSettings) (ISourceSettings, error)
	GetIDFromURL(url string) (string, error)
}
