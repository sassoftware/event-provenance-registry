// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package server

import "net/http"

func (s *Server) CreateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement me
		panic("implement me!")
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
