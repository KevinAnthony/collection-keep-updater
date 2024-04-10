package ctxu

import (
	"context"

	"github.com/atye/wikitable2json/pkg/client"
	"github.com/kevinanthony/gorps/v2/encoder"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/spf13/viper"
)

const (
	configLoaderKey ContextKey = "config_loader_ctx_key"
	httpKey         ContextKey = "http_ctx_key"
	wikiKey         ContextKey = "wiki_getter_ctx_key"
)

func GetConfigReader(cmd types.ICommand) types.IConfig {
	ctx := cmd.Context()

	value := ctx.Value(configLoaderKey)
	if cfg, ok := value.(types.IConfig); ok {
		return cfg
	}

	v := viper.New()
	ctx = context.WithValue(ctx, configLoaderKey, v)
	cmd.SetContext(ctx)

	return v
}

func GetHttpClient(cmd types.ICommand) http.Client {
	ctx := cmd.Context()

	value := ctx.Value(httpKey)
	if client, ok := value.(http.Client); ok {
		return client
	}

	client := http.NewClient(http.NewNativeClient(), encoder.NewFactory())
	ctx = context.WithValue(ctx, httpKey, client)
	cmd.SetContext(ctx)

	return client
}

func GetWikiGetter(cmd types.ICommand) client.TableGetter {
	ctx := cmd.Context()

	value := ctx.Value(wikiKey)
	if wikiGetter, ok := value.(client.TableGetter); ok {
		return wikiGetter
	}

	wikiGetter := client.NewTableGetter("noside")
	ctx = context.WithValue(ctx, wikiKey, wikiGetter)
	cmd.SetContext(ctx)

	return wikiGetter
}
