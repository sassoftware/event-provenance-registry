// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/message"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/graph-gophers/graphql-go"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

func (s *Server) CreateReceiver() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := &storage.EventReceiver{}
		err := json.NewDecoder(r.Body).Decode(rec)
		if err != nil {
			msg := err.Error()
			fmt.Println(msg)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, msg)
			return
		}

		// Check that the schema is valid.
		if rec.Schema.String() == "" {
			msg := "schema is required"
			fmt.Println(msg)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, msg)
			return
		}

		loader := gojsonschema.NewStringLoader(rec.Schema.String())
		_, err = gojsonschema.NewSchema(loader)
		if err != nil {
			msg := err.Error()
			fmt.Println(msg)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, msg)
			return
		}

		newRec, err := storage.CreateEventReceiver(s.DBConnector.Client, *rec)
		if err != nil {
			msg := err.Error()
			fmt.Println(msg)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, msg)
			return
		}

		s.kafkaCfg.MsgChannel <- message.Message{Data: message.Data{
			EventReceivers: []*storage.EventReceiver{newRec},
		}}

		// TODO: standardize responses
		render.JSON(w, r, newRec.ID)
	}
}

func (s *Server) GetReceiverByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		rec, err := storage.FindEventReceiver(s.DBConnector.Client, graphql.ID(id))
		handleGetResponse(w, r, rec, err)
	}
}
