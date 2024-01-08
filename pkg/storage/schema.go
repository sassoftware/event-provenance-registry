// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"encoding/json"
	"io"

	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
	yaml "gopkg.in/yaml.v3"
)

// Event type represents an event with various properties and a relationship to an event receiver.
type Event struct {
	ID          graphql.ID `json:"id" gorm:"type:varchar(255);primary_key;not null"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Version     string     `json:"version" gorm:"type:varchar(255);not null"`
	Release     string     `json:"release" gorm:"type:varchar(255);not null"`
	PlatformID  string     `json:"platform_id" gorm:"type:varchar(255);not null"`
	Package     string     `json:"package" gorm:"type:varchar(255);not null"`
	Description string     `json:"description" gorm:"type:varchar(255);not null"`
	Payload     types.JSON `json:"payload" gorm:"not null"`

	Success   bool       `json:"success" gorm:"not null"`
	CreatedAt types.Time `json:"created_at" gorm:"type:timestamptz; not null; default:CURRENT_TIMESTAMP"`

	EventReceiverID graphql.ID `json:"event_receiver_id" gorm:"type:varchar(255);not null"`
	EventReceiver   EventReceiver
}

// EventReceiver type represents an event receiver with various properties such as ID, name, type, version, etc...
type EventReceiver struct {
	ID          graphql.ID `json:"id" gorm:"type:varchar(255);primary_key;not null"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Type        string     `json:"type" gorm:"type:varchar(255);not null"`
	Version     string     `json:"version" gorm:"type:varchar(255);not null"`
	Description string     `json:"description" gorm:"type:varchar(255);not null"`

	Schema      types.JSON `json:"schema" gorm:"not null"`
	Fingerprint string     `json:"fingerprint" gorm:"type:varchar(255);not null"`
	CreatedAt   types.Time `json:"created_at" gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
}

// EventReceiverGroup represents a group of event receivers with various properties.
type EventReceiverGroup struct {
	ID          graphql.ID `json:"id" gorm:"type:varchar(255);primary_key;not null"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Type        string     `json:"type" gorm:"type:varchar(255);not null"`
	Version     string     `json:"version" gorm:"type:varchar(255);not null"`
	Description string     `json:"description" gorm:"type:varchar(255);not null"`
	Enabled     bool       `json:"enabled" gorm:"not null"`

	EventReceiverIDs []graphql.ID `json:"event_receiver_ids" gorm:"-"`

	CreatedAt types.Time `json:"created_at" gorm:"type:timestamptz; not null; default:CURRENT_TIMESTAMP"`
	UpdatedAt types.Time `json:"updated_at" gorm:"type:timestamptz; not null; default:CURRENT_TIMESTAMP"`
}

// EventReceiverGroupToEventReceiver represents the relationship between an
type EventReceiverGroupToEventReceiver struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`

	EventReceiverID graphql.ID `json:"event_receiver_id" gorm:"type:varchar(255);not null"`
	EventReceiver   EventReceiver

	EventReceiverGroup   EventReceiverGroup
	EventReceiverGroupID graphql.ID `json:"event_receiver_group_id" gorm:"type:varchar(255);not null"`
}

// ToJSON() function is a method defined on the `Event` struct. It converts an instance of the
// `Event` struct to a JSON string representation.
func (e *Event) ToJSON() (string, error) {
	content, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ToYAML() converts struct to a YAML string representation.
func (e *Event) ToYAML() (string, error) {
	content, err := yaml.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// EventFromJSON reads JSON data from a reader and returns an Event object.
func EventFromJSON(reader io.Reader) (*Event, error) {
	e := &Event{}
	err := json.NewDecoder(reader).Decode(e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// ToJSON() converts an instance of a JSON string representation.
func (e *EventReceiver) ToJSON() (string, error) {
	content, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ToYAML() converts struct to a YAML string representation.
func (e *EventReceiver) ToYAML() (string, error) {
	content, err := yaml.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// EventReceiverFromJSON reads JSON data from a reader and returns an EventReceiver object.
func EventReceiverFromJSON(reader io.Reader) (*EventReceiver, error) {
	e := &EventReceiver{}
	err := json.NewDecoder(reader).Decode(e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// ToJSON() converts struct to a JSON string representation.
func (e *EventReceiverGroup) ToJSON() (string, error) {
	content, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ToYAML() converts struct to a YAML string representation.
func (e *EventReceiverGroup) ToYAML() (string, error) {
	content, err := yaml.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// EventReceiverGroupFromJSON reads JSON data from a reader and returns an EventReceiver object.
func EventReceiverGroupFromJSON(reader io.Reader) (*EventReceiverGroup, error) {
	e := &EventReceiverGroup{}
	err := json.NewDecoder(reader).Decode(e)
	if err != nil {
		return nil, err
	}

	return e, nil
}
