package ctxu

import (
	"context"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/spf13/viper"
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
