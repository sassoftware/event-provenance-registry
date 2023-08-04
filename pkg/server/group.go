// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/models"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

func (s *Server) CreateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group := &models.GroupInput{}
		err := json.NewDecoder(r.Body).Decode(group)
		if err != nil {
			fmt.Println(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		newGroup, err := storage.CreateEventReceiverGroup(s.DBConnector.Client, group.EventReceiverIDs,
			group.EventReceiverGroup)
		if err != nil {
			fmt.Println(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		//TODO: write to message bus

		render.JSON(w, r, newGroup.ID)
	}
}

func (s *Server) GetGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement me
		panic("implement me!")
	}
}

func (s *Server) GetGroupByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement me
		panic("implement me!")
	}
}
