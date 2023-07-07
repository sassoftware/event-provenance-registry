package models

import (
	"time"
)

type Event struct {
	ID          string `json:"ID" gorm:"type:varchar(255);primary_key"`
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Version     string `json:"version" gorm:"type:varchar(255);not null"`
	Release     string `json:"release" gorm:"type:varchar(255);not null"`
	PlatformID  string `json:"platform_id" gorm:"type:varchar(255);not null"`
	Package     string `json:"package" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:varchar(255);not null"`
	Payload     string `json:"payload" gorm:"type:blob;not null"`

	Success   bool       `json:"success" gorm:"not null"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`

	EventReceiverID string         `json:"event_receiver_id"`
	EventReceiver   *EventReceiver `json:"event_receiver"`
}

type EventInput struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Release         string `json:"release"`
	PlatformID      string `json:"platformID"`
	Package         string `json:"package"`
	Description     string `json:"description"`
	Payload         string `json:"payload"`
	EventReceiverID string `json:"event_receiver_id"`
	Success         bool   `json:"success"`
}

type EventReceiverGroup struct {
	ID          string     `json:"ID" gorm:"type:varchar(255);primary_key"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Type        string     `json:"type" gorm:"type:varchar(255);not null"`
	Version     string     `json:"version" gorm:"type:varchar(255);not null"`
	Description string     `json:"description" gorm:"type:varchar(255);not null"`
	Enabled     bool       `json:"enabled" gorm:"not null"`
	CreatedAt   *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`

	EventReceiversID []string         `json:"event_receiver_ids"`
	EventReceivers   []*EventReceiver `json:"event_receivers"`
}

type EventReceiver struct {
	ID          string     `json:"ID" gorm:"type:varchar(255);primary_key"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Type        string     `json:"type" gorm:"type:varchar(255);not null"`
	Version     string     `json:"version" gorm:"type:varchar(255);not null"`
	Description string     `json:"description" gorm:"type:varchar(255);not null"`
	Schema      string     `json:"schema" gorm:"type:blob;not null"`
	Fingerprint string     `json:"fingerprint" gorm:"type:varchar(255);not null"`
	CreatedAt   *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
}

type EventReceiverGroupInput struct {
	Name             string   `json:"name"`
	Type             string   `json:"type"`
	Version          string   `json:"version"`
	Description      string   `json:"description"`
	Enabled          bool     `json:"enabled"`
	EventReceiverIds []string `json:"event_receiver_ids"`
}

type EventReceiverInput struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
}
