// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/sassoftware/event-provenance-registry/pkg/client"
	"github.com/sassoftware/event-provenance-registry/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger = utils.MustGetLogger("client", "cmd.root")

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "epr-cli",
	Short: "Event Provenance Registry (EPR) CLI",
	Long: `The Event Provenance Registry (EPR) CLI is a command 
	line for querying, creating, and modifying events, event-receivers, 
	and event-receiver-groups.`,
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

func preRun(cmd *cobra.Command, _ []string) error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("EPR")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}
	return nil
}

func run(_ *cobra.Command, _ []string) error {
	logger.V(1).Info("debug enabled")
	// TODO probably need some better input validation
	host := viper.GetString("host")
	port := viper.GetString("port")
	logger.Info("Host: %s", host)
	logger.Info("Port: %s", port)
	_, err := client.New(host + `:` + port)
	return err
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(filepath.Join(xdg.ConfigHome, "epr"))
		for _, dir := range xdg.ConfigDirs {
			viper.AddConfigPath(filepath.Join(dir, "epr"))
		}
		viper.SetConfigName("epr")
	}
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.MergeInConfig(); err == nil {
		fmt.Println("Merged config file:", viper.ConfigFileUsed())
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// create two new flags, one for host and one for port
	rootCmd.Flags().String("host", "localhost", "host to listen on")
	rootCmd.Flags().String("port", "8042", "port to listen on")

	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/epr/epr.yaml)")
	rootCmd.Flags().Bool("debug", false, "Enable debugging statements")
}
