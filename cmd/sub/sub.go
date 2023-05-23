// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package sub

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.sas.com/polaris/polaris/lib/utils"
)

// getLogger for logging
func getLogger() *logr.Logger {
	zerologr.NameFieldName = "logger"
	zerologr.NameSeparator = "/"
	zerologr.SetMaxV(1)

	zl := zerolog.New(os.Stdout)
	zl = zl.With().Timestamp().Logger()
	logger := zerologr.New(&zl)
	return &logger
}

// subCmd represents the subCmd command
var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "Sub CLI Command ",
	Long: `Sub is a CLI for Foo
          use it for foo`,
	PreRunE: preRun,
	RunE:    run,
}

func preRun(cmd *cobra.Command, args []string) error {
	viper.SetEnvPrefix("GENERIC")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}
	debug := viper.GetBool("debug")
	envDebug := utils.GetEnv("GENERIC_OTHER_DEBUG", "false")
	if strings.ToLower(envDebug) == `true` {
		debug = true
	}
	fmt.Print(debug)

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	logger := getLogger()
	logger.V(1).Info("If you can read this debug is on")
	logger.Info("This is a sub command")
	name := viper.GetString("name")
	version := viper.GetString("version")
	release := viper.GetString("release")
	platform := viper.GetString("platform-id")
	pkg := viper.GetString("package")
	test := viper.GetBool("test")
	logger.Info("Name", name)
	logger.Info("Version", version)
	logger.Info("Release", release)
	logger.Info("Platform ID", platform)
	logger.Info("Package", pkg)
	logger.Info("Test", test)
	if len(args) == 0 {
		return fmt.Errorf("a path to a Test File is required")
	}
	for _, arg := range args {
		logger.V(1).Info("ARG", arg)
	}
	return nil
}

// NewSubCmd func takes no input and returns *cobra.Command
// creates a new cmdline
func NewSubCmd() *cobra.Command {
	subCmd.Flags().String("name", "", "Name of the thing")
	subCmd.Flags().String("version", "", "Version of the thing")
	subCmd.Flags().String("release", "", "Release of the thing")
	subCmd.Flags().String("platform-id", "", "PlatformID of the thing")
	subCmd.Flags().String("package", "", "Package of the thing")
	subCmd.Flags().Bool("test", false, "Boolean value of test")
	return subCmd
}
