// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// GroupInput rest representation of the data for a storage.EventReceiverGroup
type GroupInput struct {
	storage.EventReceiverGroup
	EventReceiverIDs []graphql.ID `json:"event_receiver_ids"`
}

func (s *Server) CreateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := s.createGroup(r)
		handleResponse(w, r, id, err)
	}
}

func (s *Server) GetGroupByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "groupID")
		logger.V(1).Info("GetGroupByID", "groupID", id)
		rec, err := storage.FindEventReceiverGroup(s.DBConnector.Client, graphql.ID(id))
		if err != nil {
			err = missingObjectError{msg: err.Error()}
		}
		handleResponse(w, r, rec, err)
	}
}

func (s *Server) SetGroupEnabled(enabled bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "groupID")
		logger.V(1).Info("set group enabled", "groupID", id, "enabled", enabled)
		err := storage.SetEventReceiverGroupEnabled(s.DBConnector.Client, graphql.ID(id), enabled)
		if err != nil {
			err = missingObjectError{msg: err.Error()}
		}
		handleResponse(w, r, id, err)
	}
}

func (s *Server) createGroup(r *http.Request) (graphql.ID, error) {
	input := &GroupInput{}
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		return "", invalidInputError{msg: err.Error()}
	}

	eventReceiverGroupInput := storage.EventReceiverGroup{
		Name:             input.Name,
		Type:             input.Type,
		Version:          input.Version,
		Description:      input.Description,
		Enabled:          true,
		EventReceiverIDs: input.EventReceiverIDs,
	}

	eventReceiverGroup, err := storage.CreateEventReceiverGroup(s.DBConnector.Client, eventReceiverGroupInput)
	if err != nil {
		return "", err
	}

	s.kafkaCfg.MsgChannel <- message.NewEventReceiverGroup(eventReceiverGroup)
	logger.V(1).Info("created", "eventReceiverGroup", eventReceiverGroup)

	return eventReceiverGroup.ID, nil
}
