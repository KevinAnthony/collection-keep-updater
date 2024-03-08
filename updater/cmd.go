package updater

import (
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/ctxu"

	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use:   "update",
		Short: "Update Libraries based on sources",
		RunE:  run,
	}

	try   bool
	write bool
)

func init() {
	cmd.PersistentFlags().BoolVarP(&try, "print-config", "p", false, "run wanted and output the results.")
	cmd.PersistentFlags().BoolVarP(&write, "write-config", "w", false, "save the configuration to the library.")

	cmd.MarkFlagsOneRequired("print-config", "write-config")
	cmd.MarkFlagsMutuallyExclusive("print-config", "write-config")
}

func GetCmd() *cobra.Command {
	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
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

	updateSvc := New(sources)

	availableBooks, err := updateSvc.GetAllAvailableBooks(cmd.Context(), cfg.Series)
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
		case try:
			if err := library.OutputWanted(cmd, wanted); err != nil {
				return err
			}
		case write:
			if err := library.SaveWanted(wanted); err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}
