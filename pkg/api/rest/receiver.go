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
	var input epr.EventReceiverInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		return "", eprErrors.InvalidInputError{Msg: err.Error()}
	}

	eventReceiver, err := epr.CreateEventReceiver(s.msgProducer, s.DBConnector, input)
	if err != nil {
		return "", err
	}
	return eventReceiver.ID, nil
}
