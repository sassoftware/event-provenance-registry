// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/api"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/config"
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
	Use:   "server",
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

func preRun(cmd *cobra.Command, _ []string) error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}
	return nil
}

func run(_ *cobra.Command, _ []string) error {
	logger.V(1).Info("debug enabled")

	cfg, err := config.New(
		config.WithServer(viper.GetString("host"), viper.GetString("port"), "", true, true),
		config.WithStorage("localhost", "postgres", "", "", "postgres", 5432, 10, 10, 10),
		// TODO: add these once kafka and auth have been turned on
		// config.WithKafka(),
		// config.WithAuth(),
	)
	if err != nil {
		return err
	}

	ctx := context.Background()
	router, err := api.InitializeAPI(ctx, cfg)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              cfg.GetSrvAddr(),
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.Serve(listener)
		if err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			}
			logger.Info("listener closed")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		logger.Info("shutting down server")
		err := server.Shutdown(context.Background())
		if err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
		logger.Info("server shut down")
	}()

	logger.Info(fmt.Sprintf("connect to http://%s/api/v1/graphql for GraphQL playground", cfg.GetSrvAddr()))

	wg.Wait()
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

	// create two new flags, one for host and one for port
	rootCmd.Flags().String("host", "localhost", "host to listen on")
	rootCmd.Flags().String("port", "8080", "port to listen on")

	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.generic.yaml)")
	rootCmd.Flags().Bool("debug", false, "Enable debugging statements")
}
