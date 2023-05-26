// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/status"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/utils"
)

var logger = utils.MustGetLogger("server", "cmd.root")

const usage = `
server - A server for accepting events, storing events, and producing messages on a message bus.
`

var cfgFile string

// GetUsage prints version and usage
func GetUsage() {
	fmt.Println(status.GetVersion())
	fmt.Printf("%s\n", usage)
}

// GetUsageErr returns an error after printing version and usage
func GetUsageErr(err error) error {
	fmt.Printf("ERROR : %s\n", err.Error())
	GetUsage()
	return fmt.Errorf("use help")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "generic",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: preRun,
	RunE:    run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
	logger.V(1).Info("If you can read this debug is on")
	logger.Info("This is the main command")
	GetUsage()
	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Print("unable to find home directory")
		}
		// Search config in home directory with name ".generic" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".generic")
	}
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.MergeInConfig(); err == nil {
		fmt.Println("Merged config file:", viper.ConfigFileUsed())
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.generic.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debugging statements")
}
