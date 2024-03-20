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
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(config.GetCmd())
	rootCmd.AddCommand(updater.GetCmd())
}

func initConfig() {
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

	_ = viper.ReadInConfig()
}

func LoadConfig(cmd types.ICommand, _ []string) error {
	httpClient := http.NewClient(http.NewNativeClient(), encoder.NewFactory())

	sources := map[types.SourceType]types.ISource{
		types.WikipediaSource: wikipedia.New(httpClient),
		types.VizSource:       viz.New(httpClient),
		types.YenSource:       yen.New(httpClient),
		types.Kodansha:        kodansha.New(httpClient),
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
