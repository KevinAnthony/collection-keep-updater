package cmd

import (
	"github.com/kevinanthony/collection-keep-updater/config"
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/updater"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "keep-u",
		Short: "Keep your book wanted library up to date",
		Long: `Keep your book collection wanted section up to date.  
Configure it with different sources and it will compare what you already have listed with what is available and generate a wanted list.`,
		PersistentPreRunE: LoadConfig,
		RunE:              Run,
	}
)

func LoadConfig(cmd *cobra.Command, args []string) error {
	var cfg types.Config
	if err := viper.Unmarshal(&cfg, config.SeriesConfigHookFunc()); err != nil {
		return err
	}

	ctxu.SetConfig(cmd, cfg)
	ctxu.SetDI(cmd, cfg)

	return nil
}

func Run(cmd *cobra.Command, args []string) error {
	cfg, err := ctxu.GetConfig(cmd)
	if err != nil {
		return err
	}

	libraries, err := ctxu.GetLibraries(cmd)
	if err != nil {
		return err
	}

	sources, err := ctxu.GetSources(cmd)
	if err != nil {
		return err
	}

	updateSvc := updater.New(sources)

	availableBooks, err := updateSvc.GetAllAvailableBooks(cmd.Context(), cfg.Series)
	if err != nil {
		return err
	}

	for _, library := range libraries {
		err := updateSvc.UpdateLibrary(cmd.Context(), library, availableBooks)
		if err != nil {
			return err
		}
	}

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(config.GetCmd())
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME/.keepu")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
}
