// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
	"gotest.tools/v3/assert"
)

func TestEvent(t *testing.T) {

	// Test the event
	event := Event{
		ID:              "01HGDYVD995K6F24SAW6GP17HZ",
		Name:            "test",
		Version:         "0.1.1",
		Release:         "20231129",
		PlatformID:      "aarch64-gnu-linux-7",
		Package:         "OCI",
		Description:     "Test Description",
		Payload:         types.JSON{JSON: `{"name": "value"}`},
		Success:         true,
		EventReceiverID: "01HGDZ1D3KPZHYADNSJC4K4BQF",
	}

	json_out, err := event.ToJSON()
	assert.NilError(t, err)
	assert.Assert(t, strings.HasPrefix(json_out, "{"))
	assert.Assert(t, strings.HasSuffix(json_out, "}"))

	fmt.Printf("%s", json_out)

	// Test the event from JSON
	e, err := EventFromJSON(strings.NewReader(json_out))
	assert.NilError(t, err)
	assert.Assert(t, e.ID == event.ID)
	assert.Assert(t, e.Name == event.Name)
}
