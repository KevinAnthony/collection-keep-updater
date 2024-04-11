package ctxu

import (
	"context"
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/pkg/errors"
)

//go:generate mockery --srcpkg=context --name=Context --structname=ContextMock --filename=context_mock.go --output . --outpkg=ctxu

type ContextKey string

const (
	configKey    ContextKey = "config_ctx_key"
	librariesKey ContextKey = "libraries_ctx_key"
	sourcesKey   ContextKey = "sources_ctx_key"
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

func SetSources(cmd types.ICommand, sources map[types.SourceType]types.ISource) {
	ctx := cmd.Context()

	ctx = context.WithValue(ctx, sourcesKey, sources)

	cmd.SetContext(ctx)
}

func SetLibraries(cmd types.ICommand, libraries map[types.LibraryType]types.ILibrary) {
	ctx := cmd.Context()

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

func GetSource(cmd types.ICommand, sourceType types.SourceType) (types.ISource, error) {
	value := cmd.Context().Value(sourcesKey)

	sources, ok := value.(map[types.SourceType]types.ISource)
	if !ok {
		return nil, errors.New("sources not found in context")
	}

	source, found := sources[sourceType]
	if !found {
		return nil, fmt.Errorf("source type %s not found in source map", sourceType)
	}

	return source, nil
}

func GetSourceSetting(cmd types.ICommand, sourceType types.SourceType) (types.ISourceConfig, error) {
	return GetSource(cmd, sourceType)
}
