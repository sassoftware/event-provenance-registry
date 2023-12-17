// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/adrg/xdg"
	"github.com/sassoftware/event-provenance-registry/pkg/api"
	"github.com/sassoftware/event-provenance-registry/pkg/config"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/sassoftware/event-provenance-registry/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var logger = utils.MustGetLogger("server", "cmd.root")

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "epr-server",
	Short: "Event Provenance Registry (EPR) server",
	Long: `The Event Provenance Registry (EPR) server is a service
	that manages and stores events and tracks event-receivers
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
	brokers := strings.Split(viper.GetString("brokers"), ",")
	topic := viper.GetString("topic")
	host := viper.GetString("host")
	port := viper.GetString("port")
	cert := viper.GetString("tls-cert")
	key := viper.GetString("tls-key")

	dburl, err := url.Parse(viper.GetString("db"))
	if err != nil {
		return err
	}

	dbhost, dbportstr, err := net.SplitHostPort(dburl.Host)
	if err != nil {
		return err
	}

	dbport, err := strconv.Atoi(dbportstr)
	if err != nil {
		return err
	}

	messageChannel := make(chan message.Message, 1)
	defer close(messageChannel)

	cfg, err := config.New(
		config.WithServer(host, port, "", true, true),
		config.WithStorage(dbhost, "postgres", "", "", "postgres", dbport, 10, 10, 10),
		config.WithKafka(false, "3.4.0", brokers, topic, messageChannel),
		// TODO: add this once auth have been turned on
		// config.WithAuth(),
	)
	if err != nil {
		return err
	}

	dbConn, err := setupDatabase(cfg.Storage)
	if err != nil {
		return err
	}

	ctx, ccancel := context.WithCancel(context.Background())
	defer ccancel()
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-interruptChan
		signal.Stop(interruptChan)
		logger.Info("received signal", "signal", sig)
		ccancel()
	}()

	errGroup, ctx := errgroup.WithContext(ctx)

	router, err := api.Initialize(ctx, dbConn, cfg)
	if err != nil {
		return err
	}
	errGroup.Go(func() error {
		<-ctx.Done()
		return cfg.Kafka.Producer.Close()
	})

	server := &http.Server{
		Addr:              cfg.GetSrvAddr(),
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}

	if cert != "" && key != "" {
		servTLSCert, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			logger.Error(err, "invalid key pair", "certFile", cert, "keyFile", key)
			return err
		}

		// Create the TLS Config with the CA pool and enable Client certificate validation
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{servTLSCert},
			MinVersion:   tls.VersionTLS13,
		}

		// Create a Server instance to listen on port 8443 with the TLS config
		server.TLSConfig = tlsConfig

		listener = tls.NewListener(listener, tlsConfig)

		logger.Info("TLS Enabled")
	} else {
		// Run insecure if certs are not provided.
		logger.Info("TLS Not Enabled")
	}

	errGroup.Go(func() error {
		err := server.Serve(listener)
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			logger.Info("listener closed")
		}
		return nil
	})

	errGroup.Go(func() error {
		<-ctx.Done()
		logger.Info("shutting down server")
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return err
		}
		logger.Info("server shut down")
		return nil
	})

	logger.Info(fmt.Sprintf("connect to http://%s/api/v1/graphql for GraphQL playground", cfg.GetSrvAddr()))

	return errGroup.Wait()
}

func setupDatabase(cfg *config.StorageConfig) (*storage.Database, error) {
	dbConn, err := storage.New(cfg.Host, cfg.User, cfg.Pass, cfg.SSLMode, cfg.Name, cfg.Port)
	if err != nil {
		return nil, err
	}
	if err := dbConn.SyncSchema(); err != nil {
		return nil, err
	}
	return dbConn, nil
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
	rootCmd.Flags().String("brokers", "localhost:9092", "broker uris separated by commas")
	rootCmd.Flags().String("topic", "epr.dev.events", "topic to produce events on")
	rootCmd.Flags().String("db", "postgres://localhost:5432", "database connection string")
	rootCmd.Flags().String("tls-cert", "", "Path to the cert for the server")
	rootCmd.Flags().String("tls-key", "", "Path to the server key")
	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/epr/epr.yaml)")
	rootCmd.Flags().Bool("debug", false, "Enable debugging statements")
}
