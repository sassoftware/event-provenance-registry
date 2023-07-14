// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"time"

	"gorm.io/datatypes"
)

type Event struct {
	ID          string         `json:"id" gorm:"type:varchar(255);primary_key"`
	Name        string         `json:"name" gorm:"type:varchar(255);not null"`
	Version     string         `json:"version" gorm:"type:varchar(255);not null"`
	Release     string         `json:"release" gorm:"type:varchar(255);not null"`
	PlatformID  string         `json:"platform_id" gorm:"type:varchar(255);not null"`
	Package     string         `json:"package" gorm:"type:varchar(255);not null"`
	Description string         `json:"description" gorm:"type:varchar(255);not null"`
	Payload     datatypes.JSON `json:"payload" gorm:"not null"`

	Success   bool      `json:"success" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz; not null; default:CURRENT_TIMESTAMP"`

	EventReceiverID string `json:"event_receiver_id" gorm:"type:varchar(255);not null"`
	EventReceiver   EventReceiver
}

type EventReceiver struct {
	ID          string `json:"id" gorm:"type:varchar(255);primary_key"`
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Type        string `json:"type" gorm:"type:varchar(255);not null"`
	Version     string `json:"version" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:varchar(255);not null"`

	Schema      datatypes.JSON `json:"schema" gorm:"not null"`
	Fingerprint string         `json:"fingerprint" gorm:"type:varchar(255);not null"`
	CreatedAt   *time.Time     `json:"created_at" gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
}

type EventReceiverGroup struct {
	ID          string `json:"id" gorm:"type:varchar(255);primary_key"`
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Type        string `json:"type" gorm:"type:varchar(255);not null"`
	Version     string `json:"version" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:varchar(255);not null"`
	Enabled     bool   `json:"enabled" gorm:"not null"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz; not null; default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz; not null; default:CURRENT_TIMESTAMP"`
}

type EventReceiverGroupToEventReceiver struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`

	EventReceiverID string `json:"event_receiver_id" gorm:"type:varchar(255);not null"`
	EventReceiver   EventReceiver

	EventReceiverGroup   EventReceiverGroup
	EventReceiverGroupID string `json:"event_receiver_group_id" gorm:"type:varchar(255);not null"`
}
