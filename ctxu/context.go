package ctxu

import (
	"context"
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/library/libib"
	"github.com/kevinanthony/collection-keep-updater/source/kodansha"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/source/yen"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/encoder"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/pkg/errors"
)

type ctxKey string

const (
	configKey    ctxKey = "config_ctx_key"
	librariesKey ctxKey = "libraries_ctx_key"
	sourcesKey   ctxKey = "sources_ctx_key"
	httpKey      ctxKey = "http_ctx_key"
)

func SetConfig(cmd types.ICommand, cfg types.Config) {
	ctx := cmd.Context()

	ctx = context.WithValue(ctx, configKey, cfg)

	cmd.SetContext(ctx)
}

func GetConfig(cmd types.ICommand) (types.Config, error) {
	value := cmd.Context().Value(configKey)
	if cfg, ok := value.(types.Config); ok {
		return cfg, nil
	}

	return types.Config{}, errors.New("configuration not found in context")
}

func SetDI(cmd types.ICommand) {
	ctx := cmd.Context()

	httpClient := http.NewClient(http.NewNativeClient(), encoder.NewFactory())

	sources := map[types.SourceType]types.ISource{
		types.WikipediaSource: wikipedia.New(httpClient),
		types.VizSource:       viz.New(httpClient),
		types.YenSource:       yen.New(httpClient),
		types.Kodansha:        kodansha.New(httpClient),
	}

	ctx = context.WithValue(ctx, sourcesKey, sources)
	ctx = context.WithValue(ctx, httpKey, httpClient)

	cmd.SetContext(ctx)
}

func SetLibSettings(cmd types.ICommand, cfg types.Config) {
	ctx := cmd.Context()

	httpClient, ok := ctx.Value(httpKey).(http.Client)
	if !ok {
		panic("library not set in context")
	}

	libraries := map[types.LibraryType]types.ILibrary{}
	for _, setting := range cfg.Libraries {
		switch setting.Name {
		case types.LibIBLibrary:
			libraries[types.LibIBLibrary] = libib.New(setting, httpClient)
		}
	}

	ctx = context.WithValue(ctx, librariesKey, libraries)

	cmd.SetContext(ctx)
}

func GetLibraries(cmd types.ICommand) (map[types.LibraryType]types.ILibrary, error) {
	value := cmd.Context().Value(librariesKey)
	if lib, ok := value.(map[types.LibraryType]types.ILibrary); ok {
		return lib, nil
	}

	return nil, errors.New("libraries not found in context")
}

func GetSources(cmd types.ICommand) (map[types.SourceType]types.ISource, error) {
	value := cmd.Context().Value(sourcesKey)
	if source, ok := value.(map[types.SourceType]types.ISource); ok {
		return source, nil
	}

	return nil, errors.New("sources not found in context")
}

func GetSourceSetting(cmd types.ICommand, sourceType types.SourceType) (types.ISourceConfig, error) {
	value := cmd.Context().Value(sourcesKey)

	sourceMap, ok := value.(map[types.SourceType]types.ISource)
	if !ok {
		return nil, errors.New("sources not found in context")
	}

	source, ok := sourceMap[sourceType]
	if !ok {
		return nil, fmt.Errorf("source type %s not found in sources map", sourceType)
	}

	return source, nil
}
