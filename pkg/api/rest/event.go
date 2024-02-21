// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/epr"
	eprErrors "github.com/sassoftware/event-provenance-registry/pkg/errors"
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
		event, err := storage.FindEventByID(s.DBConnector.Client, graphql.ID(id))
		if err != nil {
			err = eprErrors.MissingObjectError{Msg: err.Error()}
		}
		handleResponse(w, r, event, err)
	}
}

func (s *Server) createEvent(r *http.Request) (graphql.ID, error) {
	var input epr.EventInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		return "", eprErrors.InvalidInputError{Msg: err.Error()}
	}

	event, err := epr.CreateEvent(s.msgProducer, s.DBConnector, input)
	if err != nil {
		return "", err
	}
	return event.ID, nil
}
