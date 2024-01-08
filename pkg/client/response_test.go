// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"gotest.tools/v3/assert"
)

func TestRespGraphQL(t *testing.T) {
	var resp RespGraphQL

	event := storage.Event{ID: "1", Name: "Event 1"}
	// Positive test case
	resp.Data.Events = []storage.Event{event}
	assert.Equal(t, "1", string(resp.Data.Events[0].ID))
	assert.Equal(t, "Event 1", resp.Data.Events[0].Name)

	// Negative test case
	resp.Errors = "Error occurred"
	assert.Equal(t, "Error occurred", resp.Errors)
}

func TestDecodeGraphQLRespFromJSON(t *testing.T) {
	// Positive test case
	jsonData := `{
  "data": {
    "events": [
      {
        "id": "01HGDYVD995K6F24SAW6GP17HZ",
        "name": "test",
        "version": "0.1.1",
        "release": "20231129",
        "platform_id": "aarch64-gnu-linux-7",
        "package": "OCI",
        "description": "Test Description",
        "payload": {
          "name": "value"
        },
        "success": true,
        "created_at": "2024-01-08T09:28:57.096398-05:00",
        "event_receiver_id": "01HGDZ1D3KPZHYADNSJC4K4BQF"
      }
    ]
  }
}`
	reader := strings.NewReader(jsonData)
	resp, err := DecodeGraphQLRespFromJSON(reader)
	assert.NilError(t, err)

	expected := storage.Event{
		ID:              "01HGDYVD995K6F24SAW6GP17HZ",
		Name:            "test",
		Version:         "0.1.1",
		Release:         "20231129",
		PlatformID:      "aarch64-gnu-linux-7",
		Package:         "OCI",
		Description:     "Test Description",
		Payload:         types.JSON{JSON: []byte(`{"name": "value"}`)},
		Success:         true,
		CreatedAt:       types.Time{},
		EventReceiverID: "01HGDZ1D3KPZHYADNSJC4K4BQF",
	}

	assert.Assert(t, resp.Data.Events[0].ID == expected.ID)

	// Negative test case - invalid JSON
	jsonData = `{"field1": "value1", "field2": "value2"`
	reader = strings.NewReader(jsonData)
	resp, err = DecodeGraphQLRespFromJSON(reader)
	assert.Assert(t, err != nil, "Expected error, but got nil")

	// Negative test case - empty JSON
	jsonData = `{}`
	reader = strings.NewReader(jsonData)
	res, err := DecodeGraphQLRespFromJSON(reader)
	assert.NilError(t, err)
	r := &RespGraphQL{}
	if !reflect.DeepEqual(res, r) {
		t.Errorf("Expected %+v, but got %+v", expected, resp)
	}
}
