package cmd

import (
	"github.com/kevinanthony/collection-keep-updater/config"
	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/di"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/updater"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "noside",
	Short: "Keep your book wanted library up to date",
	Long: `Keep your book collection wanted section up to date.  
Configure it with different sources and it will compare what you already have listed with what is available and generate a wanted list.`,
	PersistentPreRunE: PreREunE,
}

func PreREunE(cmd *cobra.Command, _ []string) error {
	viperConfig := ctxu.GetConfigReader(cmd)
	factory := di.GetDIFactory(cmd)

	if err := factory.Config(cmd, viperConfig); err != nil {
		return err
	}

	if err := factory.Sources(cmd); err != nil {
		return err
	}

	return ctxu.SetLibraries(cmd)
}

func init() {
	rootCmd.AddCommand(config.GetCmd())
	rootCmd.AddCommand(updater.GetCmd())
}

func GetRootCmd() types.ICommand {
	return rootCmd
}
