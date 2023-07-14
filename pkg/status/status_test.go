// Copyright © 2019, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx/types"

	"gitlab.sas.com/async-event-infrastructure/server/pkg/config"

	"gotest.tools/v3/assert"
)

func TestTheStatus(t *testing.T) {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Println(`{"data":{"alive":true}}`)
	})
	mux.HandleFunc("/healthz/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Println(`{"data":{"ready":true}}`)
	})

	svr := httptest.NewTLSServer(mux)

	defer svr.Close()
	// Test NewStatus
	topic := "re.server.topics.test.one"
	brokers := "re.server.kafka:9092,re.server.kafka:9093,re.server.kafka:9094"

	addr := svr.Listener.Addr().String()
	bits := strings.Split(addr, ":")
	if len(bits) != 2 {
		t.Fatal("could not parse test server address")
	}
	port := bits[1]
	cfg := config.New("localhost", port)

	// Create config data for application
	cfg.Debug = true
	cfg.DB = &config.DBConfig{
		Host: "postgres",
		User: "postgres",
		Name: "server",
		Port: 5432,
		Pass: "asdasd",
	}
	cfg.Kafka = &config.KafkaConfig{
		Peers: strings.Split(brokers, ","),
		Topic: topic,
	}

	var start time.Time
	status := New(start, cfg, true)

	assert.Equal(t, status.MetaData.(*Metadata).Database.Port, 5432, "Unexpected Port")
	assert.Equal(t, status.Service.Name, "server", "Unexpected Port")
}

func TestRequestResponseBody(t *testing.T) {
	// Setup http get request mocking
	mux := http.NewServeMux()
	mux.HandleFunc("/badStatus", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	mux.HandleFunc("/happypath", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, `{"data":{"alive":true}}`)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Test RequestResponseBody
	resp := types.JSONText(RequestResponseBody(ts.URL + "/happypath"))
	assert.Equal(t, strings.TrimSpace(resp.String()), `{"data":{"alive":true}}`)

	resp = RequestResponseBody("fakeurl")
	assert.Equal(t, resp.String(), `{"error" : "failed to perform Get request from 'fakeurl'"}`)

	resp = RequestResponseBody(ts.URL + "/badStatus")
	assert.Equal(t, resp.String(), `{"error" : bad status: 403 at '`+ts.URL+`/badStatus'"}`)
}
