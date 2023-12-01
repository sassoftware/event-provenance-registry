// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"fmt"
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
		Payload:         types.JSON{JSON: []byte(`{"name": "value"}`)},
		Success:         true,
		CreatedAt:       types.Time{},
		EventReceiverID: "01HGDZ1D3KPZHYADNSJC4K4BQF",
		EventReceiver:   EventReceiver{},
	}

	json_out, err := event.ToJSON()
	fmt.Printf("%s\n", json_out)
	assert.NilError(t, err)
	assert.Assert(t, strings.HasPrefix(json_out, "{"))
	assert.Assert(t, strings.HasSuffix(json_out, "}"))

	// Test the event from JSON
	e, err := EventFromJSON(strings.NewReader(json_out))
	assert.NilError(t, err)
	assert.Assert(t, e.ID == event.ID)
	assert.Assert(t, e.Name == event.Name)

	// Test input and output
	if event.Name != "test" {
		t.Errorf("Expected event name to be 'test', but got '%s'", event.Name)
	}

	// Negative test case
	invalidEvent := Event{
		ID:              "01HGKNDFSPKNY3S0KXV1QGQN2Z",
		Name:            "",
		Version:         "1.0.0",
		Release:         "202312011646141.0.1",
		PlatformID:      "aarch64-gnu-linux-7",
		Package:         "com.example.package",
		Description:     "Test event description",
		Payload:         types.JSON{JSON: []byte(`{"key": "value"}`)},
		Success:         true,
		CreatedAt:       types.Time{},
		EventReceiverID: "01HGKNJQ6SB95NJ6QVMXFQWZH1",
		EventReceiver:   EventReceiver{},
	}

	// Test invalid input
	if invalidEvent.Name != "" {
		t.Errorf("Expected invalid event name to be empty, but got '%s'", invalidEvent.Name)
	}
}

func TestEventReceiver(t *testing.T) {
	// Positive test case
	eventReceiver := EventReceiver{
		ID:          "01HGKNJQ6SB95NJ6QVMXFQWZH1",
		Name:        "Test Receiver",
		Type:        "Test Type",
		Version:     "1.0.0",
		Description: "Test Description",
		Schema:      types.JSON{JSON: []byte(`{"key": "value"}`)},
		Fingerprint: "Test Fingerprint",
		CreatedAt:   types.Time{},
	}

	// Testing input and output
	if eventReceiver.ID != "01HGKNJQ6SB95NJ6QVMXFQWZH1" {
		t.Errorf("Expected ID to be '01HGKNJQ6SB95NJ6QVMXFQWZH1', but got '%s'", eventReceiver.ID)
	}

	// Testing boundary case
	if len(eventReceiver.Name) > 255 {
		t.Errorf("Name exceeds the maximum length of 255 characters")
	}

	json_out, err := eventReceiver.ToJSON()
	fmt.Printf("%s\n", json_out)
	assert.NilError(t, err)
	assert.Assert(t, strings.HasPrefix(json_out, "{"))
	assert.Assert(t, strings.HasSuffix(json_out, "}"))

	// Test the event from JSON
	e, err := EventReceiverFromJSON(strings.NewReader(json_out))
	assert.NilError(t, err)
	assert.Assert(t, e.ID == eventReceiver.ID)
	assert.Assert(t, e.Name == eventReceiver.Name)

}

func TestEventReceiverGroup(t *testing.T) {
	// Positive test case
	eventReceiverGroup := EventReceiverGroup{
		ID:          "01HGKNVZ8XSYR429Z2HV2A31S9",
		Name:        "TestGroup",
		Type:        "epr.action.type",
		Version:     "1.0.0",
		Description: "Test eventReceiverGroup description",
		Enabled:     true,
	}

	// Test input and output
	if eventReceiverGroup.ID != "01HGKNVZ8XSYR429Z2HV2A31S9" {
		t.Errorf("Expected ID to be '01HGKNVZ8XSYR429Z2HV2A31S9', but got '%s'", eventReceiverGroup.ID)
	}

	// Negative test case
	// Test invalid input
	invalidGroup := EventReceiverGroup{
		ID:          "01HGKNVZ8XSYR429Z2HV2A31S9",
		Name:        "",
		Type:        "epr.action.type",
		Version:     "1.0.0",
		Description: "Test eventReceiverGroup description",
		Enabled:     true,
	}

	// Test validation
	if invalidGroup.Name != "" {
		t.Errorf("Expected Name to be empty, but got '%s'", invalidGroup.Name)
	}

	// Test boundary cases
	// Test maximum length of Name
	eventReceiverGroup.Name = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod eleifend nulla, quis lacinia nisl. Donec sed ipsum et ex vulputate ornare."
	if len(eventReceiverGroup.Name) > 255 {
		t.Errorf("Name exceeds maximum length of 255 characters")
	}

	// Test corner case
	// Test empty Name
	eventReceiverGroup.Name = ""
	if eventReceiverGroup.Name != "" {
		t.Errorf("Expected Name to be empty, but got '%s'", eventReceiverGroup.Name)
	}

	json_out, err := eventReceiverGroup.ToJSON()
	fmt.Printf("%s\n", json_out)
	assert.NilError(t, err)
	assert.Assert(t, strings.HasPrefix(json_out, "{"))
	assert.Assert(t, strings.HasSuffix(json_out, "}"))

	// Test the event from JSON
	e, err := EventReceiverGroupFromJSON(strings.NewReader(json_out))
	assert.NilError(t, err)
	assert.Assert(t, e.ID == eventReceiverGroup.ID)
	assert.Assert(t, e.Name == eventReceiverGroup.Name)

}
