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

func (s *Server) CreateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := &storage.Event{}
		err := json.NewDecoder(r.Body).Decode(e)
		if err != nil {
			handleResponse(w, r, nil, invalidInputError{msg: err.Error()})
			return
		}

		event, err := storage.CreateEvent(s.DBConnector.Client, *e)
		if err != nil {
			handleResponse(w, r, nil, invalidInputError{msg: err.Error()})
			return
		}

		s.kafkaCfg.MsgChannel <- message.NewEvent(event)
		logger.V(1).Info("created", "event", event)
		handleResponse(w, r, event.ID, nil)
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
