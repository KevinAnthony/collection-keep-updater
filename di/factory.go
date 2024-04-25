package di

import (
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/library/libib"
	"github.com/kevinanthony/collection-keep-updater/source/kodansha"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/source/yen"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"
)

const (
	depFactoryKey ctxu.ContextKey = "dep_factory_ctx_key"
)

//go:generate mockery --name=IDepFactory --structname=IDepFactoryMock --filename=factory_mock.go --inpackage
type IDepFactory interface {
	Config(cmd types.ICommand, icfg types.IConfig) error
	Libraries(cmd types.ICommand) error
}

type depFactory struct{}

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

	sources := map[types.SourceType]types.ISource{
		types.WikipediaSource: wikipedia.New(cmd),
		types.VizSource:       viz.New(cmd),
		types.YenSource:       yen.New(cmd),
		types.Kodansha:        kodansha.New(cmd),
	}

	ctxu.SetSources(cmd, sources)

	seriesSlice := icfg.Get("series").([]any)
	for _, data := range seriesSlice {
		series, err := utils.NewSeriesFromMap(cmd, data)
		if err != nil {
			return err
		}

		cfg.Series = append(cfg.Series, series)
	}

	libSlice := icfg.Get("libraries").([]any)
	for _, data := range libSlice {
		library, err := utils.NewLibraryFromMap(cmd, data)
		if err != nil {
			return err
		}

		cfg.Libraries = append(cfg.Libraries, library)
	}

	ctxu.SetConfig(cmd, cfg)

	return nil
}

func (f depFactory) Libraries(cmd types.ICommand) error {
	cfg, err := ctxu.GetConfig(cmd)
	if err != nil {
		return err
	}

	libraries := map[types.LibraryType]types.ILibrary{}
	for _, setting := range cfg.Libraries {
		switch setting.Name {
		case types.LibIBLibrary:
			libraries[types.LibIBLibrary] = libib.New(cmd, setting)
		}
	}

	ctxu.SetLibraries(cmd, libraries)

	return nil
}
