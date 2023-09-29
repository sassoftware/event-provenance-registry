// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"time"

	"gitlab.sas.com/async-event-infrastructure/server/pkg/message"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/utils"
)

var logger = utils.MustGetLogger("server", "config.config")

var (
	// Version of the server
	Version = "dev"
	// Commit hash of the server
	Commit = "dirty"
)

// Config contains application data for the gatekeeper application
type Config struct {
	StartTime time.Time      `json:"start_time"`
	Server    *ServerConfig  `json:"server"`
	Storage   *StorageConfig `json:"storage"`
	Kafka     *KafkaConfig   `json:"kafka"`
	Auth      *AuthConfig    `json:"-"`
}

type ServerConfig struct {
	Debug       bool   `json:"debug"`
	VerboseAPI  bool   `json:"verbose"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	ResourceDir string `json:"resources"`
}

// GetSrvAddr returns a string HOST:PORT
func (c *Config) GetSrvAddr() string {
	return c.Server.Host + ":" + c.Server.Port
}

// DBConfig holds config information about the database.
type StorageConfig struct {
	Name            string `json:"name"`
	Host            string `json:"-"`
	Port            int    `json:"-"`
	User            string `json:"-"`
	Pass            string `json:"-"`
	SSLMode         string `json:"-"`
	MaxConnections  int    `json:"db_max_connections"`
	IdleConnections int    `json:"db_idle_connections"`
	ConnectionLife  int    `json:"db_connection_max_life"`
}

// AuthConfig holds config data for authentication.
type AuthConfig struct {
	ClientID       string
	TrustedIssuers []string
}

// KafkaConfig holds config information about Kafka
type KafkaConfig struct {
	TLS        bool                 `json:"tls"`
	Version    string               `json:"version"`
	Topic      string               `json:"topic"`
	Peers      []string             `json:"peers"`
	Producer   message.Producer     `json:"-"`
	MsgChannel chan message.Message `json:"-"`
}

// LogConfigInfo Dumps most of the config info to the log.
func (c *Config) LogInfo() {
	logger.Info("Host: " + c.Server.Host)
	logger.Info("Port: " + c.Server.Port)
	logger.Info("Storage Host: " + c.Storage.Host)
	logger.Info(fmt.Sprintf("Kafka Peers: %v", c.Kafka.Peers))
	logger.Info("Kafka Version: " + c.Kafka.Version)
	logger.Info(fmt.Sprintf("Kafka TLS: %v", c.Kafka.TLS))
	logger.Info("Kafka Topic: ", c.Kafka.Topic)
	logger.Info(fmt.Sprintf("Debug: %v", c.Server.Debug))
	logger.Info(fmt.Sprintf("Verbose API: %v", c.Server.VerboseAPI))
}

// Options is a function that takes a config and returns an error
type Options func(*Config) error

// New returns a new config configured with the given options
func New(opts ...Options) (*Config, error) {
	config := &Config{}
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, err
		}
	}
	config.StartTime = time.Now()
	return config, nil
}

// WithStorage returns an option that sets the storage config
func WithStorage(host, user, pass, sslMode, name string, port, maxConnections, idleConnections, connectionLife int) Options {
	return func(cfg *Config) error {
		cfg.Storage = &StorageConfig{
			Host:            host,
			Port:            port,
			User:            user,
			Pass:            pass,
			SSLMode:         sslMode,
			Name:            name,
			MaxConnections:  maxConnections,
			IdleConnections: idleConnections,
			ConnectionLife:  connectionLife,
		}
		return nil
	}
}

// WithServer returns an option that sets the server config
func WithServer(host, port, resourceDir string, debug, verbose bool) Options {
	return func(cfg *Config) error {
		cfg.Server = &ServerConfig{
			Host:        host,
			Port:        port,
			ResourceDir: resourceDir,
			Debug:       debug,
			VerboseAPI:  verbose,
		}
		return nil
	}
}

// WithKafka returns an option that sets the kafka config
func WithKafka(tls bool, version string, peers []string, topic string, channel chan message.Message) Options {
	return func(cfg *Config) error {
		cfg.Kafka = &KafkaConfig{
			TLS:        tls,
			Version:    version,
			Peers:      peers,
			Topic:      topic,
			MsgChannel: channel,
		}
		return nil
	}
}

// WithAuth returns an option that sets the auth config
func WithAuth(clientID string, trustedIssuers []string) Options {
	return func(cfg *Config) error {
		cfg.Auth = &AuthConfig{
			ClientID:       clientID,
			TrustedIssuers: trustedIssuers,
		}
		return nil
	}
}
