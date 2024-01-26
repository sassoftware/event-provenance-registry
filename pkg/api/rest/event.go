// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

func (s *Server) CreateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := s.createEvent(r)
		handleResponse(w, r, id, err)
	}
}

func (s *Server) GetEventByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "eventID")
		event, err := storage.FindEvent(s.DBConnector.Client, graphql.ID(id))
		if err != nil {
			err = missingObjectError{msg: err.Error()}
		}
		handleResponse(w, r, event, err)
	}
}

func (s *Server) createEvent(r *http.Request) (graphql.ID, error) {
	e := &storage.Event{}
	err := json.NewDecoder(r.Body).Decode(e)
	if err != nil {
		return "", invalidInputError{msg: err.Error()}
	}

	event, err := storage.CreateEvent(s.DBConnector.Client, *e)
	if err != nil {
		return "", err
	}

	slog.Info("created", "event", event)
	s.msgProducer.Async(message.NewEvent(*event))

	eventReceiverGroups, err := storage.FindTriggeredEventReceiverGroups(s.DBConnector.Client, *event)
	if err != nil {
		slog.Error("error finding triggered event receiver groups", "error", err, "input", e)
		return "", err
	}

	for _, eventReceiverGroup := range eventReceiverGroups {
		s.msgProducer.Async(message.NewEventReceiverGroupComplete(*event, eventReceiverGroup))
	}

	slog.Info("created", "event", event)
	return event.ID, nil
}
