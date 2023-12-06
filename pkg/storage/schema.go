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
// @property ID - The ID property is of type graphql.ID and represents the unique identifier of the
// event. It is marked as the primary key in the database and cannot be null.
// @property {string} Name - The Name property represents the name of the event.
// @property {string} Version - The "Version" property represents the version of the event. It is a
// string type and is not nullable.
// @property {string} Release - The "Release" property represents the release version of the event. It
// is a string type and is not nullable.
// @property {string} PlatformID - PlatformID is a string field that represents the ID of the platform
// associated with the event.
// @property {string} Package - The "Package" property in the Event struct represents the package name
// of the event. It is of type string and is used to identify the package to which the event belongs.
// @property {string} Description - The `Description` property is a string that represents a
// description or summary of the event. It provides additional information about the event and its
// purpose.
// @property Payload - The `Payload` property is of type `types.JSON` and is used to store JSON data.
// It is marked as `not null`, meaning it must have a value.
// @property {bool} Success - The "Success" property is a boolean value that indicates whether the
// event was successful or not. It is used to track the success status of the event.
// @property CreatedAt - CreatedAt is a field that represents the timestamp when the event was created.
// It is of type `types.Time` and is stored in the database as `timestamptz` (timestamp with time
// zone). The default value for this field is set to the current timestamp when a new event is created
// @property EventReceiverID - The `EventReceiverID` property is a unique identifier for the event
// receiver associated with the event. It is of type `graphql.ID` and is stored as a string in the
// database.
// @property {EventReceiver} EventReceiver - The `EventReceiver` property is a reference to the
// `EventReceiver` struct. It represents the receiver of the event.
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

// EventReceiver type represents an event receiver with various properties such as ID, name, type, version,
// description, schema, fingerprint, and creation timestamp.
// @property ID - The ID property is of type graphql.ID and represents the unique identifier of the
// EventReceiver. It is marked as the primary key in the database and cannot be null.
// @property {string} Name - The "Name" property represents the name of the event receiver. It is a
// string type and is not nullable.
// @property {string} Type - The "Type" property in the EventReceiver struct represents the type of the
// event receiver. It is a string that describes the type of event receiver, such as "webhook",
// "email", "sms", etc.
// @property {string} Version - The "Version" property represents the version of the event receiver. It
// is a string that indicates the version of the event receiver.
// @property {string} Description - The "Description" property is a string that represents a brief
// description or summary of the event receiver. It provides additional information about the purpose
// or functionality of the event receiver.
// @property Schema - The "Schema" property is of type "types.JSON" and represents the JSON schema for
// the event receiver. It is used to define the structure and validation rules for the data that the
// event receiver expects to receive.
// @property {string} Fingerprint - The "Fingerprint" property is a string that represents a unique
// identifier or hash value for the event receiver. It is used to ensure the integrity and uniqueness
// of the event receiver.
// @property CreatedAt - CreatedAt is a property of the EventReceiver struct that represents the
// timestamp when the event receiver was created. It is of type types.Time and is stored in the
// database as timestamptz (timestamp with time zone). The default value for CreatedAt is set to the
// current timestamp when a new event
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
// @property ID - The ID property is of type graphql.ID and represents the unique identifier of the
// EventReceiverGroup. It is marked as the primary key in the database and cannot be null.
// @property {string} Name - The Name property represents the name of the EventReceiverGroup.
// @property {string} Type - The "Type" property in the EventReceiverGroup struct represents the type
// of the event receiver group. It is a string that describes the category or purpose of the group.
// @property {string} Version - The "Version" property in the EventReceiverGroup struct represents the
// version of the event receiver group. It is a string type and is used to indicate the version of the
// group.
// @property {string} Description - The "Description" property is a string that represents a brief
// description or summary of the EventReceiverGroup. It provides additional information about the
// purpose or functionality of the group.
// @property {bool} Enabled - The "Enabled" property is a boolean value that indicates whether the
// EventReceiverGroup is enabled or not. If it is set to true, it means the group is enabled and can
// receive events. If it is set to false, it means the group is disabled and will not receive any
// events.
// @property {[]graphql.ID} EventReceiverIDs - EventReceiverIDs is a slice of graphql.ID that
// represents the IDs of the event receivers associated with the EventReceiverGroup. This property is
// not stored in the database as it is marked with `gorm:"-"`, indicating that it is not a database
// column.
// @property CreatedAt - CreatedAt is a property of type `types.Time` that represents the timestamp
// when the EventReceiverGroup was created. It is stored in the database as a timestamptz (timestamp
// with time zone) and has a default value of the current timestamp.
// @property UpdatedAt - UpdatedAt is a property of the EventReceiverGroup struct. It represents the
// timestamp of when the EventReceiverGroup was last updated. The value of UpdatedAt is of type
// types.Time, which is a custom type that likely represents a timestamp. The gorm tags indicate that
// the UpdatedAt field should be mapped
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
// `EventReceiverGroup` and an `EventReceiver` in Go.
// @property {int} ID - The ID property is an integer field that serves as the primary key for the
// EventReceiverGroupToEventReceiver struct. It is annotated with `json:"id"` to specify the JSON key
// for this field when marshaling and unmarshaling JSON data. It is also annotated with
// `gorm:"primaryKey;
// @property EventReceiverID - The `EventReceiverID` property is of type `graphql.ID` and represents
// the ID of the event receiver associated with this event receiver group to event receiver mapping.
// @property {EventReceiver} EventReceiver - The `EventReceiver` property is a reference to an
// `EventReceiver` object. It represents the event receiver associated with this event receiver group
// to event receiver relationship.
// @property {EventReceiverGroup} EventReceiverGroup - The `EventReceiverGroup` property represents the
// event receiver group associated with the `EventReceiverGroupToEventReceiver` object. It is of type
// `EventReceiverGroup`.
// @property EventReceiverGroupID - The EventReceiverGroupID property is an identifier for the
// EventReceiverGroup that this EventReceiverGroupToEventReceiver belongs to.
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

// ToYAML() function is a method defined on the `Event` struct. It converts an instance of the
// `Event` struct to a YAML string representation.
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

// ToJSON() function is a method defined on the `EventReceiver` struct. It converts an instance of the
// `EventReceiver` struct to a JSON string representation.
func (e *EventReceiver) ToJSON() (string, error) {
	content, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ToYAML() function is a method defined on the `EventReceiver` struct. It converts an instance
// of the `EventReceiver` struct to a YAML string representation.
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

// ToJSON() function is a method defined on the `EventReceiverGroup` struct. It converts an
// instance of the `EventReceiverGroup` struct to a JSON string representation.
func (e *EventReceiverGroup) ToJSON() (string, error) {
	content, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ToYAML() function is a method defined on the `EventReceiverGroup` struct. It converts an
// instance of the `EventReceiverGroup` struct to a YAML string representation.
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
