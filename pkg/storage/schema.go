// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
)

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

type EventReceiverGroupToEventReceiver struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`

	EventReceiverID graphql.ID `json:"event_receiver_id" gorm:"type:varchar(255);not null"`
	EventReceiver   EventReceiver

	EventReceiverGroup   EventReceiverGroup
	EventReceiverGroupID graphql.ID `json:"event_receiver_group_id" gorm:"type:varchar(255);not null"`
}
