// Copyright © 2019, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/sassoftware/event-provenance-registry/pkg/config"
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

	// create config object
	cfg, err := config.New(
		config.WithServer("localhost", port, "/resources", true, true),
		config.WithStorage("postgres", "user", "pass", "ssl", "postgres", 5432, 10, 10, 10),
		config.WithKafka(true, "2.6", strings.Split(brokers, ","), topic, nil),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create create status config from config object
	status := New(cfg)

	assert.Equal(t, status.MetaData.(*Metadata).Storage.Port, 5432, "Unexpected Port")
	assert.Equal(t, status.MetaData.(*Metadata).Storage.Host, "postgres", "Unexpected Host")
	assert.Equal(t, status.MetaData.(*Metadata).Storage.Name, "postgres", "Unexpected Name")
	assert.Equal(t, status.MetaData.(*Metadata).Storage.User, "user", "Unexpected User")
	assert.Equal(t, status.MetaData.(*Metadata).Storage.Pass, "pass", "Unexpected Pass")
	assert.Equal(t, status.MetaData.(*Metadata).Storage.SSLMode, "ssl", "Unexpected SSLMode")
	assert.Equal(t, status.MetaData.(*Metadata).Storage.MaxConnections, 10, "Unexpected MaxConnections")
	assert.Equal(t, status.MetaData.(*Metadata).Storage.IdleConnections, 10, "Unexpected IdleConnections")
	assert.Equal(t, status.MetaData.(*Metadata).Storage.ConnectionLife, 10, "Unexpected ConnectionLife")
	assert.DeepEqual(t, status.MetaData.(*Metadata).Kafka.Peers, strings.Split(brokers, ","))
	assert.Equal(t, status.MetaData.(*Metadata).Kafka.Topic, topic, "Unexpected Topic")
	assert.Equal(t, status.MetaData.(*Metadata).Kafka.Version, "2.6", "Unexpected Version")
	assert.Equal(t, status.MetaData.(*Metadata).Verbose, true, "Unexpected Verbose")
	assert.Equal(t, status.MetaData.(*Metadata).Resources, "/resources", "Unexpected Resources")

	assert.Equal(t, status.Service.Name, "server", "Unexpected name")
	assert.Equal(t, status.Service.Version, "dev", "Unexpected version")
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
	assert.Equal(t, resp.String(), `{"error" : "bad status: 403 at '`+ts.URL+`/badStatus'"}`)
}
