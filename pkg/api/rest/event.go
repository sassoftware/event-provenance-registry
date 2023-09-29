// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/json"
	"fmt"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/message"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/graph-gophers/graphql-go"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

func (s *Server) CreateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event := &storage.Event{}
		err := json.NewDecoder(r.Body).Decode(event)
		if err != nil {
			fmt.Println(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		newEvent, err := storage.CreateEvent(s.DBConnector.Client, *event)
		if err != nil {
			fmt.Println(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		s.kafkaCfg.MsgChannel <- message.Message{Data: message.Data{
			Events: []*storage.Event{newEvent},
		}}
		render.JSON(w, r, newEvent.ID)
	}
}

func (s *Server) GetEventByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		event, err := storage.FindEvent(s.DBConnector.Client, graphql.ID(id))
		handleGetResponse(w, r, event, err)
	}
}
