package db

import "time"

// Base is an adaptation of the basic gorm model
type Base struct {
	CreatedAt time.Time `gorm:"type:timestamptz; not null; default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamptz; not null; default:CURRENT_TIMESTAMP"`
}

type Event struct {
	ID          string `gorm:"uniqueIndex:event_pk;type:varchar(255);primary_key"`
	Name        string `gorm:"uniqueIndex:event_pk;type:varchar(255);not null"`
	Version     string `gorm:"uniqueIndex:event_pk;type:varchar(255);not null"`
	Release     string `gorm:"uniqueIndex:event_pk;type:varchar(255);not null"`
	PlatformID  string `gorm:"uniqueIndex:event_pk;type:varchar(255);not null"`
	Package     string `gorm:"uniqueIndex:event_pk;type:varchar(255);not null"`
	Description string `gorm:"type:varchar(255);not null"`
	Payload     string `gorm:"type:varchar(255);not null"`
	Success     bool   `gorm:"not null"`

	EventReceiverID string `gorm:"type:varchar(255);not null"`
	EventReceiver   EventReceiver

	Base
}

type EventReceiver struct {
	ID          string `gorm:"uniqueIndex:event_receiver_pk;type:varchar(255);primary_key"`
	Name        string `gorm:"uniqueIndex:event_receiver_pk;type:varchar(255);not null"`
	Type        string `gorm:"uniqueIndex:event_receiver_pk;type:varchar(255);not null"`
	Version     string `gorm:"uniqueIndex:event_receiver_pk;type:varchar(255);not null"`
	Description string `gorm:"type:varchar(255);not null"`
	Enabled     bool   `gorm:"not null"`

	Base
}

type EventReceiverGroup struct {
	ID          string `gorm:"uniqueIndex:event_receiver_group_pk;type:varchar(255);primary_key"`
	Name        string `gorm:"uniqueIndex:event_receiver_group_pk;type:varchar(255);not null"`
	Type        string `gorm:"uniqueIndex:event_receiver_group_pk;type:varchar(255);not null"`
	Version     string `gorm:"uniqueIndex:event_receiver_group_pk;type:varchar(255);not null"`
	Description string `gorm:"type:varchar(255);not null"`
	Enabled     bool   `gorm:"not null"`

	Base
}

type EventReceiverGroupToEventReceiver struct {
	ID int `gorm:"primaryKey;autoIncrement"`

	EventReceiverID string `gorm:"uniqueIndex:receiver_group_link;type:varchar(255);not null"`
	EventReceiver   EventReceiver

	EventReceiverGroup   EventReceiverGroup
	EventReceiverGroupID string `gorm:"uniqueIndex:receiver_group_link;type:varchar(255);not null"`
}
