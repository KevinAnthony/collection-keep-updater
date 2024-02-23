package ctxu

import (
	"context"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type ctxKey string

const (
	configKey ctxKey = "config_ctx_key"
)

func SetConfigCtx(cmd *cobra.Command, cfg types.Config) {
	ctx := cmd.Context()

	ctx = context.WithValue(ctx, configKey, cfg)

	cmd.SetContext(ctx)
}

func GetConfigCtx(cmd *cobra.Command) (types.Config, error) {
	value := cmd.Context().Value(configKey)
	if cfg, ok := value.(types.Config); ok {
		return cfg, nil
	}

	return types.Config{}, errors.New("configuration not found in context")
}
