// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/epr"
	eprErrors "github.com/sassoftware/event-provenance-registry/pkg/errors"
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
		slog.Info("GetGroupByID", "groupID", id)
		rec, err := storage.FindEventReceiverGroupByID(s.DBConnector.Client, graphql.ID(id))
		handleResponse(w, r, rec, err)
	}
}

func (s *Server) UpdateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "groupID")
		slog.Info("update group", "groupID", id)

		var patch struct {
			Enabled *bool `json:"enabled,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
			err = eprErrors.InvalidInputError{Msg: err.Error()}
			handleResponse(w, r, id, err)
			return
		}

		var err error
		if patch.Enabled != nil {
			slog.Info("set group enabled", "groupID", id, "enabled", patch.Enabled)
			err = storage.SetEventReceiverGroupEnabled(s.DBConnector.Client, graphql.ID(id), *patch.Enabled)
		}
		handleResponse(w, r, id, err)
	}
}

func (s *Server) createGroup(r *http.Request) (graphql.ID, error) {
	var input epr.EventReceiverGroupInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		return "", eprErrors.InvalidInputError{Msg: err.Error()}
	}

	eventReceiverGroup, err := epr.CreateEventReceiverGroup(s.msgProducer, s.DBConnector, input)
	if err != nil {
		return "", err
	}
	return eventReceiverGroup.ID, nil
}
