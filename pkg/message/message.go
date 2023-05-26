// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package message

import (
	"encoding/json"
	"io"

	"gitlab.sas.com/async-event-infrastructure/server/pkg/models"
	yaml "gopkg.in/yaml.v3"
)

// Message is the struct for kafka message events it contains information from the receipt that created the event
// Adheres to https://github.com/cloudevents/spec 1.0.1
type Message struct {
	Success     bool   `json:"success"`     // Extension to Cloud Events Spec
	ID          string `json:"id"`          // Cloud Events Spec 1.0.1
	Specversion string `json:"specversion"` // Cloud Events Spec 1.0.1
	Type        string `json:"type"`        // Cloud Events Spec 1.0.1
	Source      string `json:"source"`      // Cloud Events Spec 1.0.1
	APIVersion  string `json:"api_version"` // Extension to Cloud Events Spec
	Name        string `json:"name"`        // Extension to Cloud Events Spec
	Version     string `json:"version"`     // Extension to Cloud Events Spec
	Release     string `json:"release"`     // Extension to Cloud Events Spec
	PlatformID  string `json:"platform_id"` // Extension to Cloud Events Spec
	Package     string `json:"package"`     // Extension to Cloud Events Spec
	Data        Data   `json:"data"`        // Cloud Events Spec 1.0.1
}

// Data contains the data that created the event
type Data struct {
	Events         []*models.Event              `json:"events"`
	EventReceivers []*models.EventReceiver      `json:"event_receivers"`
	EventGroups    []*models.EventReceiverGroup `json:"event_groups"`
}

// ToJSON converts a Events struct to JSON
func (m *Message) ToJSON() (string, error) {
	content, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ToYAML converts a Events struct to YAML
func (m *Message) ToYAML() (string, error) {
	content, err := yaml.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// New returns an Message
func New() *Message {
	return &Message{
		Specversion: CloudEventsSpec,
		APIVersion:  APIv1,
	}
}

// DecodeFromJSON returns an Event from JSON
func DecodeFromJSON(reader io.Reader) (*Message, error) {
	message := &Message{}
	err := json.NewDecoder(reader).Decode(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// MsgInfo is a base type for encoding messages
type MsgInfo interface {
	Length() int
	Encode() ([]byte, error)
}

type messageInfo struct {
	msgType interface{}
	encoded []byte
	err     error
	Topic   string `json:"-"`
}

func (m *messageInfo) ensureEncoded() {
	if m.encoded == nil && m.err == nil {
		m.encoded, m.err = json.Marshal(m.msgType)
	}
}

// Length returns the length of the encoded message
func (m *messageInfo) Length() int {
	m.ensureEncoded()
	return len(m.encoded)
}

// Encode encodes the message
func (m *messageInfo) Encode() ([]byte, error) {
	m.ensureEncoded()
	return m.encoded, m.err
}
