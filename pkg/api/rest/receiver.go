// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/graph-gophers/graphql-go"
	eprErrors "github.com/sassoftware/event-provenance-registry/pkg/errors"
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
		slog.Info("getting receiver", "id", id)
		eventReceiver, err := storage.FindEventReceiver(s.DBConnector.Client, graphql.ID(id))
		handleResponse(w, r, eventReceiver, err)
	}
}

func (s *Server) createReceiver(r *http.Request) (graphql.ID, error) {
	rec := &storage.EventReceiver{}
	err := json.NewDecoder(r.Body).Decode(rec)
	if err != nil {
		return "", eprErrors.InvalidInputError{Msg: err.Error()}
	}

	// Check that the schema is valid.
	if rec.Schema.String() == "" {
		return "", eprErrors.InvalidInputError{Msg: "schema is required"}
	}

	loader := gojsonschema.NewStringLoader(rec.Schema.String())
	_, err = gojsonschema.NewSchema(loader)
	if err != nil {
		return "", eprErrors.InvalidInputError{Msg: err.Error()}
	}

	eventReceiver, err := storage.CreateEventReceiver(s.DBConnector.Client, *rec)
	if err != nil {
		return "", err
	}

	s.msgProducer.Async(message.NewEventReceiver(eventReceiver))
	slog.Info("created", "eventReceiver", eventReceiver)
	return eventReceiver.ID, nil
}
