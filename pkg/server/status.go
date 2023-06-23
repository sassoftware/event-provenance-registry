// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/go-chi/render"
	"net/http"
)

func (s *Server) CheckLiveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, `{"alive"":true}`)
	}
}

func (s *Server) CheckReadiness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, `{"ready":true}`)
	}
}

func (s *Server) CheckStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, `{"ready":true}`)
	}
}
