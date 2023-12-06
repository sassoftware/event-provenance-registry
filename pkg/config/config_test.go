// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewConfig(t *testing.T) {
	cfg, err := New(
		WithServer("localhost", "8080", "/resources", true, true),
		WithStorage("clash.london.com", "joe", "brixton", "disable", "postgres", 5432, 10, 10, 10),
		WithKafka(true, "2.6", []string{"kafka.svc.cluster.local:9092"}, "server.events", nil),
		WithAuth("01HGX8QDVTMSXXQHNV9AH7X8QQ", []string{"foo", "bar"}),
	)

	if err != nil {
		t.Fatal(err)
	}
	cfg.LogInfo()

	assert.Assert(t, cfg.Server.Host == "localhost", "Expected host to be 'localhost', got %s", cfg.Server.Host)
	assert.Assert(t, cfg.Server.Port == "8080", "Expected port to be '8080', got %s", cfg.Server.Port)
	assert.Assert(t, cfg.Storage.Host == "clash.london.com", "Expected host to be 'clash.london.com', got %s", cfg.Storage.Host)
	assert.Assert(t, cfg.Storage.Port == 5432, "Expected port to be 5432, got %d", cfg.Storage.Port)
	assert.Assert(t, cfg.Storage.User == "joe", "Expected user to be 'joe', got %s", cfg.Storage.User)
	assert.Assert(t, cfg.Storage.Pass == "brixton", "Expected pass to be 'pass', got %s", cfg.Storage.Pass)
	assert.Assert(t, cfg.Storage.SSLMode == "disable", "Expected sslMode to be 'disable', got %s", cfg.Storage.SSLMode)
	assert.Assert(t, cfg.Storage.Name == "postgres", "Expected name to be 'postgres', got %s", cfg.Storage.Name)
	assert.Assert(t, cfg.Storage.MaxConnections == 10, "Expected maxConnections to be 10, got %d", cfg.Storage.MaxConnections)
	assert.Assert(t, cfg.Storage.IdleConnections == 10, "Expected idleConnections to be 10, got %d", cfg.Storage.IdleConnections)
	assert.Assert(t, cfg.Storage.ConnectionLife == 10, "Expected connectionLife to be 10, got %d", cfg.Storage.ConnectionLife)
	assert.Assert(t, cfg.Kafka.Version == "2.6", "Expected kafka version to be '2.6', got %s", cfg.Kafka.Version)
	assert.Assert(t, cfg.Auth.ClientID == "01HGX8QDVTMSXXQHNV9AH7X8QQ", "Expected auth client id to be '01HGX8QDVTMSXXQHNV9AH7X8QQ', got %s", cfg.Auth.ClientID)
}
