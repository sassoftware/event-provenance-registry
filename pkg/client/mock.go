// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// make sure the mock implements the interface
var _ Contract = &MockEPR{}

// MockEPR a fake implementation of the gatekeeper for unit testing. Change the EventReceiver/EventReceiverGroup/Event
// fields to determine what gets returned.
type MockEPR struct {
	Error              error
	Event              *storage.Event
	EventReceiver      *storage.EventReceiver
	EventReceiverGroup *storage.EventReceiverGroup
}

// CreateEventReceiver fakes gate creation
func (m MockEPR) CreateEventReceiver(_ *storage.EventReceiver) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}

	panic("implement me")
}

// CreateEvent fakes receipt creation
func (m MockEPR) CreateEvent(_ *storage.Event) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}

	panic("implement me")
}

// CreateEventReceiverGroup fakes stage creation
func (m MockEPR) CreateEventReceiverGroup(_ *storage.EventReceiverGroup) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}

	panic("implement me")
}

// ModifyEventReceiverGroup fakes stage modification
func (m MockEPR) ModifyEventReceiverGroup(_ *storage.EventReceiverGroup) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}

	panic("implement me")
}

// ModifyEvent fakes receipt modification
func (m MockEPR) ModifyEvent(_ *storage.Event) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}

	panic("implement me")
}

// CheckReadiness fake readiness check
func (m MockEPR) CheckReadiness() (bool, error) {
	if m.Error != nil {
		return false, m.Error
	}

	panic("implement me")
}

// CheckLiveness fake liveness check
func (m MockEPR) CheckLiveness() (bool, error) {
	if m.Error != nil {
		return false, m.Error
	}

	panic("implement me")
}

// GetEndpoint fakes getting an endpoint
func (m MockEPR) GetEndpoint(_ string) string {
	panic("implement me")
}

// QueryEventReceiverByID fakes getting a gate
func (m MockEPR) QueryEventReceiverByID(_ string, _ []string) (*storage.EventReceiver, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	panic("implement me")
}

// QueryEventByID fakes getting a receipt
func (m MockEPR) QueryEventByID(_ string, _ []string) (*storage.Event, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	return m.Event, nil
}

// QueryEventReceiverGroupByID fakes getting a stage
func (m MockEPR) QueryEventReceiverGroupByID(_ string, _ []string) (*storage.EventReceiverGroup, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	panic("implement me")
}

// SearchEventReceiverGroups searching for a stage
func (m MockEPR) SearchEventReceiverGroups(_ map[string]interface{}, _ []string) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}

	panic("implement me")
}

// SearchEventReceiverGroupsObj searching for a stage
func (m MockEPR) SearchEventReceiverGroupsObj(_ map[string]interface{}, _ []string) ([]storage.EventReceiverGroup, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	panic("implement me")
}

// SearchEventReceivers fakes searching for a EventReceiver
func (m MockEPR) SearchEventReceivers(_ map[string]interface{}, _ []string) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}

	panic("implement me")
}

// SearchEventReceiversObj fakes searching for a EventReceiver
func (m MockEPR) SearchEventReceiversObj(_ map[string]interface{}, _ []string) ([]storage.EventReceiver, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	panic("implement me")
}

// SearchEvents fakes searching for a Event
func (m MockEPR) SearchEvents(_ map[string]interface{}, _ []string) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}
	panic("implement me")
}

// SearchEventsObj fakes searching for a Event
func (m MockEPR) SearchEventsObj(_ map[string]interface{}, _ []string) ([]storage.Event, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	panic("implement me")
}

// Search fakes a search
func (m MockEPR) Search(_ string, _ string, _ map[string]interface{},
	_ []string) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}
	panic("implement me")
}

// IDSearch fakes a search
func (m MockEPR) IDSearch(_ map[string]interface{}) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}
	panic("implement me")
}
