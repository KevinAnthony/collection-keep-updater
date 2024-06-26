package updater

import (
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/spf13/cobra"
)

const (
	printF = "print-config"
	writeF = "write-config"
)

var cmd = &cobra.Command{
	Use:   "update",
	Short: "Update Libraries based on sources",
	RunE:  types.CmdRun(Run),
}

func init() {
	cmd.PersistentFlags().BoolP(printF, "p", false, "run wanted and output the results.")
	cmd.PersistentFlags().BoolP(writeF, "w", false, "save the configuration to the library.")

	cmd.MarkFlagsOneRequired("print-config", "write-config")
	cmd.MarkFlagsMutuallyExclusive("print-config", "write-config")
}

func GetCmd() *cobra.Command {
	return cmd
}

func Run(cmd types.ICommand) error {
	cfg, err := ctxu.GetConfig(cmd)
	if err != nil {
		return err
	}

	libraries, err := ctxu.GetLibraries(cmd)
	if err != nil {
		return err
	}

	updateSvc := GetUpdater(cmd)

	availableBooks, err := updateSvc.GetAllAvailableBooks(cmd, cfg.Series)
	if err != nil {
		return err
	}

	for _, library := range libraries {
		wanted, err := updateSvc.UpdateLibrary(cmd.Context(), library, availableBooks)
		if err != nil {
			return err
		}

		switch {
		case len(wanted) == 0:
			fmt.Println("No New Wanted books")

			continue
		case utils.GetFlagBool(cmd, printF):
			wanted.Print(cmd)
		case utils.GetFlagBool(cmd, writeF):
			if err := library.SaveWanted(cmd, wanted); err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}
