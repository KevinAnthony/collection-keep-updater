package cmd

import (
	"github.com/kevinanthony/collection-keep-updater/config"
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/source/kodansha"
	"github.com/kevinanthony/collection-keep-updater/source/viz"
	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/source/yen"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/updater"
	"github.com/kevinanthony/gorps/v2/encoder"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/atye/wikitable2json/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "noside",
	Short: "Keep your book wanted library up to date",
	Long: `Keep your book collection wanted section up to date.  
Configure it with different sources and it will compare what you already have listed with what is available and generate a wanted list.`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		viperConfig := viper.New()
		httpClient := http.NewClient(http.NewNativeClient(), encoder.NewFactory())
		wikiGetter := client.NewTableGetter("noside")

		if err := LoadConfig(cmd, viperConfig); err != nil {
			return err
		}

		return LoadDI(cmd, httpClient, wikiGetter)
	},
}

func init() {
	rootCmd.AddCommand(config.GetCmd())
	rootCmd.AddCommand(updater.GetCmd())
}

func LoadDI(cmd types.ICommand, httpClient http.Client, wikiGetter client.TableGetter) error {
	sources := map[types.SourceType]types.ISource{}

	if wiki, err := wikipedia.New(httpClient, wikiGetter); err != nil {
		return err
	} else {
		sources[types.WikipediaSource] = wiki
	}

	if vizSource, err := viz.New(httpClient); err != nil {
		return err
	} else {
		sources[types.VizSource] = vizSource
	}

	if yenSource, err := yen.New(httpClient); err != nil {
		return err
	} else {
		sources[types.YenSource] = yenSource
	}

	if kodanshaSource, err := kodansha.New(httpClient); err != nil {
		return err
	} else {
		sources[types.Kodansha] = kodanshaSource
	}

	ctxu.SetDI(cmd, httpClient, sources)

	return ctxu.SetLibraries(cmd)
}

func LoadConfig(
	cmd types.ICommand,
	icfg types.IConfig,
) error {
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

	ctxu.SetConfig(cmd, icfg, cfg)

	return nil
}

func GetRootCmd() types.ICommand {
	return rootCmd
}
