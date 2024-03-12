package types

import (
	"context"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ ICommand = (*cobra.Command)(nil)

//go:generate mockery --name=ICommand --structname=ICommandMock --filename=cmd_mock.go --inpackage
type ICommand interface {
	Context() context.Context
	SetContext(ctx context.Context)
	OutOrStdout() io.Writer
	Flag(name string) (flag *pflag.Flag)
	PersistentFlags() *pflag.FlagSet
	Flags() *pflag.FlagSet
}

func CmdRunE(f func(cmd ICommand, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return f(cmd, args)
	}
}

func CmdArgs(f func(cmd ICommand, args []string) error) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		return f(cmd, args)
	}
}

func CmdPersistentPreRunE(f func(cmd ICommand, _ []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return f(cmd, args)
	}
}
