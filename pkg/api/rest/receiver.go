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
	"github.com/xeipuuv/gojsonschema"
)

func (s *Server) CreateReceiver() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := &storage.EventReceiver{}
		err := json.NewDecoder(r.Body).Decode(rec)
		if err != nil {
			handleResponse(w, r, nil, invalidInputError{msg: err.Error()})
			return
		}

		// Check that the schema is valid.
		if rec.Schema.String() == "" {
			err := invalidInputError{msg: "schema is required"}
			handleResponse(w, r, nil, err)
			return
		}

		loader := gojsonschema.NewStringLoader(rec.Schema.String())
		_, err = gojsonschema.NewSchema(loader)
		if err != nil {
			handleResponse(w, r, nil, invalidInputError{msg: err.Error()})
			return
		}

		eventReceiver, err := storage.CreateEventReceiver(s.DBConnector.Client, *rec)
		if err != nil {
			handleResponse(w, r, nil, invalidInputError{msg: err.Error()})
			return
		}

		s.kafkaCfg.MsgChannel <- message.NewEventReceiver(eventReceiver)
		logger.V(1).Info("created", "eventReceiver", eventReceiver)
		handleResponse(w, r, eventReceiver.ID, nil)
	}
}

func (s *Server) GetReceiverByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "receiverID")
		logger.V(1).Info("GetReceiverByID", "receiverID", id)
		eventReceiver, err := storage.FindEventReceiver(s.DBConnector.Client, graphql.ID(id))
		if err != nil {
			err = missingObjectError{msg: err.Error()}
		}
		handleResponse(w, r, eventReceiver, err)
	}
}
