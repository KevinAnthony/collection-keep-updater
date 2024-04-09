package di

import (
	"context"

	"github.com/kevinanthony/collection-keep-updater/config"
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/source/kodansha"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/source/yen"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/atye/wikitable2json/pkg/client"
)

const (
	depFactoryKey ctxu.ContextKey = "dep_factory_ctx_key"
)

//go:generate mockery --name=IDepFactory --structname=IDepFactoryMock --filename=di_mock.go --inpackage
type IDepFactory interface {
	Sources(cmd types.ICommand, httpClient http.Client, wikiGetter client.TableGetter) error
	Config(cmd types.ICommand, icfg types.IConfig) error
}

type depFactory struct{}

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

func (depFactory) Sources(cmd types.ICommand, httpClient http.Client, wikiGetter client.TableGetter) error {
	vizSource, err := viz.New(httpClient)
	if err != nil {
		return err
	}

	wikiSource, err := wikipedia.New(httpClient, wikiGetter)
	if err != nil {
		return err
	}

	yenSource, err := yen.New(httpClient)
	if err != nil {
		return err
	}

	kodanshaSource, err := kodansha.New(httpClient)
	if err != nil {
		return err
	}

	sources := map[types.SourceType]types.ISource{
		types.WikipediaSource: wikiSource,
		types.VizSource:       vizSource,
		types.YenSource:       yenSource,
		types.Kodansha:        kodanshaSource,
	}

	ctxu.SetDI(cmd, httpClient, sources)

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
