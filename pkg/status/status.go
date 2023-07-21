// Copyright © 2019, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx/types"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/config"
)

// Health reports the links to the health monitors
type Health struct {
	Readiness types.JSONText `json:"readiness" yaml:"readiness"`
	Liveness  types.JSONText `json:"liveness" yaml:"liveness"`
}

// Status reports the characteristics of the Service
type Status struct {
	Service   *Version    `json:"service" yaml:"service"`
	Uptime    string      `json:"uptime" yaml:"uptime"`
	Debug     bool        `json:"debug" yaml:"debug"`
	Health    *Health     `json:"health" yaml:"health"`
	Info      *Info       `json:"info" yaml:"info"`
	StartTime time.Time   `json:"start_time" yaml:"start_time"`
	Host      string      `json:"host" yaml:"host"`
	Port      string      `json:"port" yaml:"port"`
	MetaData  interface{} `json:"metadata" yaml:"metadata"`
}

// Metadata struct for info about service-specific fields
type Metadata struct {
	Verbose   bool                `json:"verbose" yaml:"verbose"`
	Resources string              `json:"resources" yaml:"resources"`
	Kafka     *config.KafkaConfig `json:"kafka" yaml:"kafka"`
	Database  *config.DBConfig    `json:"database" yaml:"database"`
}

// GetStatus returns a string of JSON equivalent
// to the given Status
func (s *Status) GetStatus() string {
	content, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return `{"error" : "failed to marshal status"}`
	}
	return string(content)
}

// NewServerStatus creates a new Status struct for use in a Status service
func NewMetadata(cfg *config.Config) *Metadata {
	return &Metadata{
		Resources: cfg.ResourceDir,
		Kafka:     cfg.Kafka,
		Database:  cfg.DB,
		Verbose:   cfg.Verbose,
	}
}

// New will return a new Status struct for service
// given a name, version, release, debug variable, and start time.
func New(start time.Time, cfg *config.Config, debug bool) *Status {
	server := cfg.GetSrvAddr()
	metadata := NewMetadata(cfg)
	stat := &Status{
		Service:   NewVersion(),
		Debug:     debug,
		StartTime: start,
		Uptime:    GetUptime(start),
		Health:    NewHealth(server),
		Info:      GetInfo(),
		MetaData:  metadata,
	}

	return stat
}

// NewHealth will get the health of the service given a server address
// using the provided address to hit that services' endpoints
// TODO use built in functions for this instead of calling endpoints
func NewHealth(server string) *Health {
	return &Health{
		Readiness: RequestResponseBody("http://" + server + "/healthz/readiness"),
		Liveness:  RequestResponseBody("http://" + server + "/healthz/liveness"),
	}
}

// RequestResponseBody will perform an http Get
// using the provided url and return the body as a []byte
// TODO use built in functions for this instead of calling endpoints
func RequestResponseBody(url string) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []byte(`{"error" : "failed formulate request for '` + url + `'"}`)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("error: %s", err)
		return []byte(`{"error" : "failed to perform Get request from '` + url + `'"}`)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte(`{"error" : bad status: ` + fmt.Sprint(resp.StatusCode) + ` at '` + url + `'"}`)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte(`{"error" : unable to read response body from ` + url + `}`)
	}
	return body
}

// GetUptime returns uptime of server
func GetUptime(start time.Time) string {
	return time.Since(start).String()
}
