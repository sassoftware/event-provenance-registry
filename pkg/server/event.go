// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
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
		// TODO: write to message bus

		render.JSON(w, r, newEvent.ID)
	}
}

func (s *Server) GetEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement me
		panic("implement me!")
	}
}

func (s *Server) GetEventByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement me
		panic("implement me!")
	}
}
