package ctxu

import (
	"context"

	"github.com/kevinanthony/collection-keep-updater/source/yen"

	"github.com/kevinanthony/collection-keep-updater/library/libib"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/encoder"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type ctxKey string

const (
	configKey    ctxKey = "config_ctx_key"
	librariesKey ctxKey = "libraries_ctx_key"
	sorucesKey   ctxKey = "sources_ctx_key"
)

func SetConfig(cmd *cobra.Command, cfg types.Config) {
	ctx := cmd.Context()

	ctx = context.WithValue(ctx, configKey, cfg)

	cmd.SetContext(ctx)
}

func GetConfig(cmd *cobra.Command) (types.Config, error) {
	value := cmd.Context().Value(configKey)
	if cfg, ok := value.(types.Config); ok {
		return cfg, nil
	}

	return types.Config{}, errors.New("configuration not found in context")
}

func SetDI(cmd *cobra.Command, cfg types.Config) {
	ctx := cmd.Context()

	httpClient := http.NewClient(http.NewNativeClient(), encoder.NewFactory())

	libraries := map[types.LibraryType]types.ILibrary{}
	for _, setting := range cfg.Libraries {
		switch setting.Name {
		case types.LibIBLibrary:
			libraries[types.LibIBLibrary] = libib.New(setting, httpClient)
		}
	}
	sources := map[types.SourceType]types.ISource{
		types.WikipediaSource: wikipedia.New(httpClient),
		types.VizSource:       viz.New(httpClient),
		types.YenSource:       yen.New(httpClient),
	}

	ctx = context.WithValue(ctx, librariesKey, libraries)
	ctx = context.WithValue(ctx, sorucesKey, sources)

	cmd.SetContext(ctx)
}

func GetLibraries(cmd *cobra.Command) (map[types.LibraryType]types.ILibrary, error) {
	value := cmd.Context().Value(librariesKey)
	if lib, ok := value.(map[types.LibraryType]types.ILibrary); ok {
		return lib, nil
	}

	return nil, errors.New("libraries not found in context")
}

func GetSources(cmd *cobra.Command) (map[types.SourceType]types.ISource, error) {
	value := cmd.Context().Value(sorucesKey)
	if source, ok := value.(map[types.SourceType]types.ISource); ok {
		return source, nil
	}

	return nil, errors.New("sources not found in context")
}
