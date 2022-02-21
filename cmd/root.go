package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/must"
	"murphysec-cli-simple/version"
	"os"
	"strings"
)

var versionFlag bool

func rootCmd() *cobra.Command {
	argsMap := map[string]bool{}
	for _, it := range os.Args {
		argsMap[it] = true
	}
	c := &cobra.Command{
		Use:               "murphysec",
		PersistentPreRunE: preRun,
		TraverseChildren:  true,
		Run: func(cmd *cobra.Command, args []string) {
			must.Must(cmd.Help())
		},
	}
	c.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "show version and exit")
	c.PersistentFlags().BoolVar(&logger.DisableLogFile, "no-log-file", false, "do not write log file")
	c.PersistentFlags().StringVar(&logger.CliLogFilePathOverride, "write-log-to", "", "specify log file path")
	c.PersistentFlags().StringVar(&logger.ConsoleLogLevelOverride, "log-level", "", "specify log level, must be silent|error|warn|info|debug")
	c.PersistentFlags().StringVar(&conf.APITokenCliOverride, "token", "", "specify API token")
	c.PersistentFlags().StringVar(&api.CliServerAddressOverride, "server", "", "specify server address")
	c.PersistentFlags().String("ide", "", "hidden")
	must.Must(c.PersistentFlags().MarkHidden("ide"))
	c.AddCommand(authCmd())
	c.AddCommand(scanCmd())
	c.AddCommand(ideaScanCmd())
	return c
}

func preRun(cmd *cobra.Command, args []string) error {
	if versionFlag {
		version.PrintVersionInfo()
		os.Exit(0)
	}
	if !utils.InStringSlice([]string{"", "warn", "error", "debug", "info", "silent"}, strings.ToLower(strings.TrimSpace(logger.ConsoleLogLevelOverride))) {
		return errors.New("Loglevel invalid, must be silent|error|warn|info|debug")
	}
	return nil
}

func Execute() error {
	return rootCmd().Execute()
}
