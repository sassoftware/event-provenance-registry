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

// GroupInput rest representation of the data for a storage.EventReceiverGroup
type GroupInput struct {
	storage.EventReceiverGroup
	EventReceiverIDs []graphql.ID `json:"event_receiver_ids"`
}

func (s *Server) CreateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &GroupInput{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			fmt.Println(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		eventRecieverGroupInput := storage.EventReceiverGroup{
			Name:             input.Name,
			Type:             input.Type,
			Version:          input.Version,
			Description:      input.Description,
			Enabled:          true,
			EventReceiverIDs: input.EventReceiverIDs,
		}

		eventReceiverGroup, err := storage.CreateEventReceiverGroup(s.DBConnector.Client, eventRecieverGroupInput)
		if err != nil {
			fmt.Println(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		s.kafkaCfg.MsgChannel <- message.Message{Data: message.Data{
			EventGroups: []*storage.EventReceiverGroup{eventReceiverGroup},
		}}
		render.JSON(w, r, eventReceiverGroup.ID)
	}
}

func (s *Server) GetGroupByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		rec, err := storage.FindEventReceiverGroup(s.DBConnector.Client, graphql.ID(id))
		handleGetResponse(w, r, rec, err)
	}
}

func (s *Server) SetGroupEnabled(enabled bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		err := storage.SetEventReceiverGroupEnabled(s.DBConnector.Client, graphql.ID(id), enabled)

		if err != nil {
			fmt.Println(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, RestResponse{Errors: []error{err}})
			return
		}
		render.JSON(w, r, RestResponse{})
	}
}
