// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg, err := New(
		WithServer("localhost", "8080", "/resources", true, true),
		WithStorage("clash.london.com", "joe", "brixton", "foo", "posgres", 5432, 10, 10, 10),
		WithKafka(true, "2.6", []string{"kafka.svc.cluster.local:9092"}, "server.events", nil, nil),
		WithAuth("ABCDEFGHIJK", []string{"foo", "bar"}),
	)

	if err != nil {
		t.Fatal(err)
	}
	cfg.LogInfo()
}
