// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func customStringCompare(query string, options []string) cmp.Comparison {
	return func() cmp.Result {
		for _, possible := range options {
			if possible == query {
				return cmp.ResultSuccess
			}
		}
		return cmp.ResultFailure(fmt.Sprintf("%q did not match available options %v", query, options))
	}
}

// Full querys examples
// {"query":"query ($erg: FindEventReceiverGroupInput!){event_receiver_groups(event_receiver_group: $erg) {id,name,type,version,description}}","variables":{"erg": {"name":"foobar","version":"1.0.0"}}}
// {"query":"query ($er: FindEventReceiverInput!){event_receivers(event_receiver: $er) {id,name,type,version,description}}","variables":{"er":{"id":"01HPW652DSJBHR5K4KCZQ97GJP"}}}
// {"query":"query ($e : FindEventInput!){events(event: $e) {id,name,version,release,platform_id,package,description,success,event_receiver_id}}","variables":{"e": {"name":"foo","version":"1.0.0"}}}
func TestNewGraphQLSearchRequest(t *testing.T) {
	operation := "events"
	fields := []string{"id", "name", "version", "release", "platform_id", "package", "description", "success", "event_receiver_id"}
	params := map[string]interface{}{
		"name":    "foo",
		"version": "1.0.0",
	}
	expected := []string{`query ($obj: FindEventInput!){events(event: $obj) {id,name,version,release,platform_id,package,description,success,event_receiver_id}}`}
	req := NewGraphQLSearchRequest(operation, params, fields)
	assert.Check(t, customStringCompare(req.Query, expected))

	operation = "event_receivers"
	fields = []string{"id", "name", "type", "version", "description"}
	expected = []string{`query ($obj: FindEventReceiverInput!){event_receivers(event_receiver: $obj) {id,name,type,version,description}}`}
	req = NewGraphQLSearchRequest(operation, params, fields)
	assert.Check(t, customStringCompare(req.Query, expected))

	operation = "event_receiver_groups"
	expected = []string{`query ($obj: FindEventReceiverGroupInput!){event_receiver_groups(event_receiver_group: $obj) {id,name,type,version,description}}`}
	req = NewGraphQLSearchRequest(operation, params, fields)
	assert.Check(t, customStringCompare(req.Query, expected))
}

// mutation ($er: CreateEventReceiverInput!){create_event_receiver(event_receiver: $er)}
func TestNewGraphQLMutationRequest(t *testing.T) {
	operation := "create_event"
	params := map[string]interface{}{
		"name":    "foo",
		"version": "1.0.0",
	}
	expected := []string{`mutation ($obj: CreateEventInput!){create_event(event: $obj)}`}
	req := NewGraphQLMutationRequest(operation, params)
	assert.Check(t, customStringCompare(req.Query, expected))

	operation = "create_event_receiver"
	expected = []string{`mutation ($obj: CreateEventReceiverInput!){create_event_receiver(event_receiver: $obj)}`}
	req = NewGraphQLMutationRequest(operation, params)
	assert.Check(t, customStringCompare(req.Query, expected))

	operation = "create_event_receiver_group"
	expected = []string{`mutation ($obj: CreateEventReceiverGroupInput!){create_event_receiver_group(event_receiver_group: $obj)}`}
	req = NewGraphQLMutationRequest(operation, params)
	assert.Check(t, customStringCompare(req.Query, expected))
}
