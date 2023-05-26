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
	Debug         bool         `json:"debug"`
	Verbose       bool         `json:"verbose"`
	Host          string       `json:"host"`
	Port          string       `json:"port"`
	ResourceDir   string       `json:"resources"`
	URI           string       `json:"uri"`
	AzureClientID string       `json:"azure_client_id"`
	StartTime     time.Time    `json:"start_time"`
	DB            *DBConfig    `json:"database"`
	Kafka         *KafkaConfig `json:"kafka"`
	Auth          *AuthConfig  `json:"-"`
}

// DBConfig holds config information about the database.
type DBConfig struct {
	Host            string `json:"-"`
	User            string `json:"-"`
	Name            string `json:"name"`
	Port            int    `json:"-"`
	Pass            string `json:"-"`
	SSLMode         string `json:"-"`
	MaxConnections  int    `json:"db_max_connections"`
	IdleConnections int    `json:"db_idle_connections"`
	ConnectionLife  int    `json:"db_connection_max_life"`
}

// AuthConfig holds config data for authentication.
type AuthConfig struct {
	// OIDC parameters
	Issuer         string
	ClientID       string
	TrustedIssuers []string
}

// KafkaConfig holds config information about Kafka
type KafkaConfig struct {
	TLS        bool                 `json:"tls"`
	Version    string               `json:"version"`
	ClientID   string               `json:"client_id"`
	Mechanism  string               `json:"mechanism"`
	Topic      string               `json:"topic"`
	Topics     []string             `json:"topics"`
	Peers      []string             `json:"peers"`
	Producer   message.Producer     `json:"-"`
	MsgChannel chan message.Message `json:"-"`
}

// GetSrvAddr returns a string HOST:PORT
func (c *Config) GetSrvAddr() string {
	srvaddr := c.Host + ":" + c.Port
	return srvaddr
}

// LogConfigInfo Dumps most of the config info to the log.
func (c *Config) LogConfigInfo() {
	logger.Info("Host: " + c.Host)
	logger.Info("Port: " + c.Port)
	logger.Info("URI: " + c.URI)
	logger.Info("Database Host: " + c.DB.Host)
	logger.Info("Auth Service URL: " + c.Auth.Issuer)
	logger.Info(fmt.Sprintf("Kafka Peers: %v", c.Kafka.Peers))
	logger.Info("Kafka Version: " + c.Kafka.Version)
	logger.Info(fmt.Sprintf("Kafka TLS: %v", c.Kafka.TLS))
	logger.Info("Kafka Topic: " + c.Kafka.Topic)
	logger.V(1).Info(fmt.Sprintf("Debug: %v", c.Debug))
	logger.V(1).Info(fmt.Sprintf("Verbose: %v", c.Verbose))
}

// NewConfig returns a new instance of Config
func NewConfig(host string, port string) *Config {
	return &Config{
		Host:      host,
		Port:      port,
		StartTime: time.Now(),
	}
}
