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

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "noside",
		Short: "Keep your book wanted library up to date",
		Long: `Keep your book collection wanted section up to date.  
Configure it with different sources and it will compare what you already have listed with what is available and generate a wanted list.`,
		PersistentPreRunE: types.CmdPersistentPreRunE(LoadConfig),
	}
)

func init() {
	rootCmd.AddCommand(config.GetCmd())
	rootCmd.AddCommand(updater.GetCmd())
}

func LoadConfig(cmd types.ICommand, _ []string) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME/.config/noside/")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	httpClient := http.NewClient(http.NewNativeClient(), encoder.NewFactory())
	wikiGetter := client.NewTableGetter("keep-updater")
	sources := map[types.SourceType]types.ISource{}

	if wiki, err := wikipedia.New(httpClient, wikiGetter); err != nil {
		return err
	} else {
		sources[types.WikipediaSource] = wiki
	}

	if viz, err := viz.New(httpClient); err != nil {
		return err
	} else {
		sources[types.VizSource] = viz
	}

	if yen, err := yen.New(httpClient); err != nil {
		return err
	} else {
		sources[types.YenSource] = yen
	}

	if kodansha, err := kodansha.New(httpClient); err != nil {
		return err
	} else {
		sources[types.Kodansha] = kodansha
	}

	ctxu.SetDI(cmd, httpClient, sources)
	var cfg types.Config

	if err := viper.Unmarshal(&cfg, config.SeriesConfigHookFunc(cmd)); err != nil {
		return err
	}

	ctxu.SetConfig(cmd, cfg)
	if err := ctxu.SetLibraries(cmd, cfg); err != nil {
		return err
	}

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
