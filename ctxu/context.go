package ctxu

import (
	"context"
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/library/libib"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/pkg/errors"
)

//go:generate mockery --srcpkg=context --name=Context --structname=ContextMock --filename=context_mock.go --output . --outpkg=ctxu

type ContextKey string

const (
	configKey    ContextKey = "config_ctx_key"
	iconfigKey   ContextKey = "i_config_ctx_key"
	librariesKey ContextKey = "libraries_ctx_key"
	sourcesKey   ContextKey = "sources_ctx_key"
	updaterKey   ContextKey = "updater_ctx_key"
	httpKey      ContextKey = "http_ctx_key"
)

func SetConfig(cmd types.ICommand, v types.IConfig, cfg types.Config) {
	ctx := cmd.Context()

	ctx = context.WithValue(ctx, configKey, cfg)
	ctx = context.WithValue(ctx, iconfigKey, v)

	cmd.SetContext(ctx)
}

func GetConfig(cmd types.ICommand) (types.Config, error) {
	value := cmd.Context().Value(configKey)
	if cfg, ok := value.(types.Config); ok {
		return cfg, nil
	}

	return types.Config{}, errors.New("configuration not found in context")
}

func GetConfigReader(cmd types.ICommand) (types.IConfig, error) {
	value := cmd.Context().Value(iconfigKey)
	if cfg, ok := value.(types.IConfig); ok {
		return cfg, nil
	}

	return nil, errors.New("configuration reader not found in context")
}

func SetDI(cmd types.ICommand, httpClient http.Client, sources map[types.SourceType]types.ISource) {
	ctx := cmd.Context()

	ctx = context.WithValue(ctx, sourcesKey, sources)
	ctx = context.WithValue(ctx, httpKey, httpClient)

	cmd.SetContext(ctx)
}

func SetLibraries(cmd types.ICommand) error {
	ctx := cmd.Context()

	cfg, err := GetConfig(cmd)
	if err != nil {
		return err
	}

	httpClient, ok := ctx.Value(httpKey).(http.Client)
	if !ok {
		return errors.New("http client not set in context")
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

	return nil
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
