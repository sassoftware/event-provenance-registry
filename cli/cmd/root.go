// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/sassoftware/event-provenance-registry/cli/cmd/event"
	"github.com/sassoftware/event-provenance-registry/cli/cmd/group"
	"github.com/sassoftware/event-provenance-registry/cli/cmd/receiver"
	"github.com/sassoftware/event-provenance-registry/cli/cmd/status"
	"github.com/sassoftware/event-provenance-registry/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
	// TODO probably need some better input validation
	url := viper.GetString("url")
	_, err := client.New(url)
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
	eventCmd := event.NewEventCmd()
	rootCmd.AddCommand(eventCmd)
	receiverCmd := receiver.NewReceiverCmd()
	rootCmd.AddCommand(receiverCmd)
	groupCmd := group.NewGroupCmd()
	rootCmd.AddCommand(groupCmd)
	statusCmd := status.NewStatusCmd()
	rootCmd.AddCommand(statusCmd)

	rootCmd.Flags().String("url", "http://localhost:8042", "EPR base url")

	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/epr/epr.yaml)")
	rootCmd.Flags().Bool("debug", false, "Enable debugging statements")
}
