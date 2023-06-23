// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type Event struct {
	ID            string `json:"ID"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	Release       string `json:"release"`
	PlatformID    string `json:"platformID"`
	Package       string `json:"package"`
	Description   string `json:"description"`
	Payload       string `json:"payload"`
	EventReceiver string `json:"event_receiver"`
	Success       bool   `json:"success"`
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

type EventReceiver struct {
	ID          string `json:"ID"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
	Fingerprint string `json:"fingerprint"`
}

type EventReceiverGroup struct {
	ID             string   `json:"ID"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	Version        string   `json:"version"`
	Description    string   `json:"description"`
	Enabled        bool     `json:"enabled"`
	EventReceivers []string `json:"event_receivers"`
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
