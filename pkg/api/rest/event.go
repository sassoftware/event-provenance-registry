// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

func (s *Server) CreateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := &storage.Event{}
		err := json.NewDecoder(r.Body).Decode(e)
		if err != nil {
			handleGetResponse(w, r, nil, err)
			return
		}

		event, err := storage.CreateEvent(s.DBConnector.Client, *e)
		if err != nil {
			handleGetResponse(w, r, nil, err)
			return
		}

		s.kafkaCfg.MsgChannel <- message.NewEvent(event)

		logger.V(1).Info("created", "event", event)

		render.JSON(w, r, event.ID)
	}
}

func (s *Server) GetEventByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "eventID")
		event, err := storage.FindEvent(s.DBConnector.Client, graphql.ID(id))
		handleGetResponse(w, r, event, err)
	}
}
