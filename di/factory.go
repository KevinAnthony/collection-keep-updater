package di

import (
	"context"

	"github.com/kevinanthony/collection-keep-updater/config"
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/library/libib"
	"github.com/kevinanthony/collection-keep-updater/source/kodansha"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/source/yen"
	"github.com/kevinanthony/collection-keep-updater/types"
)

const (
	depFactoryKey ctxu.ContextKey = "dep_factory_ctx_key"
)

//go:generate mockery --name=IDepFactory --structname=IDepFactoryMock --filename=factory_mock.go --inpackage
type IDepFactory interface {
	Sources(cmd types.ICommand) error
	Config(cmd types.ICommand, icfg types.IConfig) error
	Libraries(cmd types.ICommand) error
}

type depFactory struct{}

func (f depFactory) Libraries(cmd types.ICommand) error {
	cfg, err := ctxu.GetConfig(cmd)
	if err != nil {
		return err
	}

	httpClient := ctxu.GetHttpClient(cmd)

	libraries := map[types.LibraryType]types.ILibrary{}
	for _, setting := range cfg.Libraries {
		switch setting.Name {
		case types.LibIBLibrary:
			libraries[types.LibIBLibrary] = libib.New(setting, httpClient)
		}
	}

	ctxu.SetLibraries(cmd, libraries)

	return nil
}

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

func (depFactory) Sources(cmd types.ICommand) error {
	sources := map[types.SourceType]types.ISource{
		types.WikipediaSource: wikipedia.New(cmd),
		types.VizSource:       viz.New(cmd),
		types.YenSource:       yen.New(cmd),
		types.Kodansha:        kodansha.New(cmd),
	}

	ctxu.SetSources(cmd, sources)

	return nil
}

func (depFactory) Config(cmd types.ICommand, icfg types.IConfig) error {
	icfg.AddConfigPath("$HOME/.config/noside/")
	icfg.AddConfigPath(".")
	icfg.SetConfigType("yaml")
	icfg.SetConfigName("config")
	icfg.AutomaticEnv()

	if err := icfg.ReadInConfig(); err != nil {
		return err
	}

	var cfg types.Config

	seriesSlice := icfg.Get("series").([]interface{})
	for _, data := range seriesSlice {
		series, err := config.GetSeries(cmd, data)
		if err != nil {
			return err
		}

		cfg.Series = append(cfg.Series, series)
	}

	libSlice := icfg.Get("libraries").([]interface{})
	for _, data := range libSlice {
		library, err := config.GetLibrary(cmd, data)
		if err != nil {
			return err
		}

		cfg.Libraries = append(cfg.Libraries, library)
	}

	ctxu.SetConfig(cmd, cfg)

	return nil
}
