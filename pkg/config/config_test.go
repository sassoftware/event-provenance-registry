// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"
)

func TestNewConfig(_ *testing.T) {
	cfg := New("http://localhost", "9420")

	// Create config data for application
	cfg.Debug = true
	cfg.Verbose = true
	cfg.ResourceDir = "/resources"
	cfg.URI = "http://server.io"
	cfg.DB = &DBConfig{
		Host:            "clash.london.com",
		User:            "joe",
		Name:            "brixton",
		Port:            5432,
		Pass:            "magnificent_seven",
		SSLMode:         "foo",
		MaxConnections:  0,
		IdleConnections: 2,
		ConnectionLife:  0,
	}
	cfg.Kafka = &KafkaConfig{
		TLS:     true,
		Version: "2.6",
		Peers:   []string{"kafka.svc.cluster.local:9092"},
		Topic:   "server.events",
	}
	cfg.Auth = &AuthConfig{
		ClientID:       "ABCDEFGHIJK",
		TrustedIssuers: []string{"foo", "bar"},
	}

	cfg.LogInfo()
}
