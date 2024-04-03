package types

import (
	"context"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	_ ICommand = (*cobra.Command)(nil)
	_ IConfig  = (*viper.Viper)(nil)
)

//go:generate mockery --name=ICommand --structname=ICommandMock --filename=cmd_mock.go --inpackage
type ICommand interface {
	Context() context.Context
	SetContext(ctx context.Context)
	OutOrStdout() io.Writer
	Flag(name string) *pflag.Flag
	PersistentFlags() *pflag.FlagSet
	Execute() error
}

//go:generate mockery --name=IConfig --structname=IConfigMock --filename=config_mock.go --inpackage
type IConfig interface {
	AddConfigPath(in string)
	SetConfigType(in string)
	SetConfigName(in string)
	AutomaticEnv()
	ReadInConfig() error
	WriteConfig() error
	Set(key string, value any)
	Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) error
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
