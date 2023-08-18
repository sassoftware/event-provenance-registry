// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
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
			fmt.Println(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		if rec.Schema.String() == "" {
			msg := "schema is required"
			fmt.Println(msg)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, msg)
			return
		}

		// TODO: validate the schema

		newRec, err := storage.CreateEventReceiver(s.DBConnector.Client, *rec)
		if err != nil {
			fmt.Println(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		// TODO: write to message bus

		// TODO: standardize responses
		render.JSON(w, r, newRec.ID)
	}
}

func (s *Server) GetReceivers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement me
		panic("implement me!")
	}
}

func (s *Server) GetReceiverByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		rec, err := storage.FindEventReceiver(s.DBConnector.Client, graphql.ID(id))
		handleReadResponse(w, r, rec, err)
	}
}
