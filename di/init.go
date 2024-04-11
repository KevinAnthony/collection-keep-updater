package di

import (
	"context"

	"github.com/kevinanthony/collection-keep-updater/types"
)

func NewDepFactory() IDepFactory {
	return depFactory{}
}

func GetDIFactory(cmd types.ICommand) IDepFactory {
	ctx := cmd.Context()

	value := ctx.Value(depFactoryKey)
	if cfg, ok := value.(IDepFactory); ok {
		return cfg
	}

	v := NewDepFactory()
	ctx = context.WithValue(ctx, depFactoryKey, v)
	cmd.SetContext(ctx)

	return v
}
