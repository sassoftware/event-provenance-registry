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
		id, err := s.createReceiver(r)
		handleResponse(w, r, id, err)
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

func (s *Server) createReceiver(r *http.Request) (graphql.ID, error) {
	rec := &storage.EventReceiver{}
	err := json.NewDecoder(r.Body).Decode(rec)
	if err != nil {
		return "", invalidInputError{msg: err.Error()}
	}

	// Check that the schema is valid.
	if rec.Schema.String() == "" {
		return "", invalidInputError{msg: "schema is required"}
	}

	loader := gojsonschema.NewStringLoader(rec.Schema.String())
	_, err = gojsonschema.NewSchema(loader)
	if err != nil {
		return "", invalidInputError{msg: err.Error()}
	}

	eventReceiver, err := storage.CreateEventReceiver(s.DBConnector.Client, *rec)
	if err != nil {
		return "", err
	}

	s.msgProducer.Async(message.NewEventReceiver(eventReceiver))
	logger.V(1).Info("created", "eventReceiver", eventReceiver)
	return eventReceiver.ID, nil
}
