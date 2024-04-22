// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package resolvers

import (
	"testing"

	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
	"gotest.tools/v3/assert"
)

var (
	id              = graphql.ID("test.id")
	name            = "test.name"
	version         = "1.0.0"
	release         = "03202024"
	platformID      = "test.platform.id"
	pkg             = "test.package"
	success         = true
	eventReceiverID = graphql.ID("test.event.receiver.id")
	typeName        = "test.type"
)

func TestEvent(t *testing.T) {
	event := EventInput{
		Name:            "test",
		Version:         "1.0.0",
		Release:         "1.0.0",
		PlatformID:      "test",
		Package:         "test",
		Description:     "test",
		Payload:         types.JSON{},
		Success:         true,
		EventReceiverID: "test",
	}

	assert.DeepEqual(t, event, EventInput{
		Name:            "test",
		Version:         "1.0.0",
		Release:         "1.0.0",
		PlatformID:      "test",
		Package:         "test",
		Description:     "test",
		Payload:         types.JSON{},
		Success:         true,
		EventReceiverID: "test",
	})
}

func TestFindEvent(t *testing.T) {
	findEvent := FindEventInput{
		ID:              &id,
		Name:            graphql.NullString{Value: &name, Set: true},
		Version:         graphql.NullString{Value: &version, Set: true},
		Release:         graphql.NullString{Value: &release, Set: true},
		PlatformID:      graphql.NullString{Value: &platformID, Set: true},
		Package:         graphql.NullString{Value: &pkg, Set: true},
		Success:         graphql.NullBool{Value: &success, Set: true},
		EventReceiverID: &eventReceiverID,
	}

	assert.DeepEqual(t, findEvent, FindEventInput{
		ID:              &id,
		Name:            graphql.NullString{Value: &name, Set: true},
		Version:         graphql.NullString{Value: &version, Set: true},
		Release:         graphql.NullString{Value: &release, Set: true},
		PlatformID:      graphql.NullString{Value: &platformID, Set: true},
		Package:         graphql.NullString{Value: &pkg, Set: true},
		Success:         graphql.NullBool{Value: &success, Set: true},
		EventReceiverID: &eventReceiverID,
	})
}

func TestFindEventToMap(t *testing.T) {
	findEvent := FindEventInput{
		ID:              &id,
		Name:            graphql.NullString{Value: &name, Set: true},
		Version:         graphql.NullString{Value: &version, Set: true},
		Release:         graphql.NullString{Value: &release, Set: true},
		PlatformID:      graphql.NullString{Value: &platformID, Set: true},
		Package:         graphql.NullString{Value: &pkg, Set: true},
		Success:         graphql.NullBool{Value: &success, Set: true},
		EventReceiverID: &eventReceiverID,
	}

	assert.DeepEqual(t, findEvent.toMap(), map[string]any{
		"id":                id,
		"name":              &name,
		"version":           &version,
		"release":           &release,
		"platform_id":       &platformID,
		"package":           &pkg,
		"success":           &success,
		"event_receiver_id": eventReceiverID,
	})
}

func TestEventReceiver(t *testing.T) {
	eventReceiver := EventReceiverInput{
		Name:        "test",
		Type:        "test",
		Version:     "1.0.0",
		Description: "test",
		Schema:      types.JSON{},
	}
	assert.DeepEqual(t, eventReceiver, EventReceiverInput{
		Name:        "test",
		Type:        "test",
		Version:     "1.0.0",
		Description: "test",
		Schema:      types.JSON{},
	})
}

func TestFindEventReceiver(t *testing.T) {
	findEventReceiver := FindEventReceiverInput{
		ID:      &id,
		Name:    graphql.NullString{Value: &name, Set: true},
		Type:    graphql.NullString{Value: &typeName, Set: true},
		Version: graphql.NullString{Value: &version, Set: true},
	}
	assert.DeepEqual(t, findEventReceiver, FindEventReceiverInput{
		ID:      &id,
		Name:    graphql.NullString{Value: &name, Set: true},
		Type:    graphql.NullString{Value: &typeName, Set: true},
		Version: graphql.NullString{Value: &version, Set: true},
	})
}

func TestFindEventReceiverToMap(t *testing.T) {
	findEventReceiver := FindEventReceiverInput{
		ID:      &id,
		Name:    graphql.NullString{Value: &name, Set: true},
		Type:    graphql.NullString{Value: &typeName, Set: true},
		Version: graphql.NullString{Value: &version, Set: true},
	}
	assert.DeepEqual(t, findEventReceiver.toMap(), map[string]any{
		"id":      id,
		"name":    &name,
		"type":    &typeName,
		"version": &version,
	})
}

func TestEventReceiverGroup(t *testing.T) {
	eventReceiverGroup := EventReceiverGroupInput{
		Name:             "test",
		Type:             "test",
		Version:          "1.0.0",
		Description:      "test",
		EventReceiverIDs: []graphql.ID{"test"},
	}
	assert.DeepEqual(t, eventReceiverGroup, EventReceiverGroupInput{
		Name:             "test",
		Type:             "test",
		Version:          "1.0.0",
		Description:      "test",
		EventReceiverIDs: []graphql.ID{"test"},
	})
}

func TestFindEventReceiverGroup(t *testing.T) {
	findEventReceiverGroup := FindEventReceiverGroupInput{
		ID:      &id,
		Name:    graphql.NullString{Value: &name, Set: true},
		Type:    graphql.NullString{Value: &typeName, Set: true},
		Version: graphql.NullString{Value: &version, Set: true},
	}
	assert.DeepEqual(t, findEventReceiverGroup, FindEventReceiverGroupInput{
		ID:      &id,
		Name:    graphql.NullString{Value: &name, Set: true},
		Type:    graphql.NullString{Value: &typeName, Set: true},
		Version: graphql.NullString{Value: &version, Set: true},
	})
}

func TestFindEventReceiverGroupToMap(t *testing.T) {
	findEventReceiverGroup := FindEventReceiverGroupInput{
		ID:      &id,
		Name:    graphql.NullString{Value: &name, Set: true},
		Type:    graphql.NullString{Value: &typeName, Set: true},
		Version: graphql.NullString{Value: &version, Set: true},
	}
	assert.DeepEqual(t, findEventReceiverGroup.toMap(), map[string]any{
		"id":      id,
		"name":    &name,
		"type":    &typeName,
		"version": &version,
	})
}
