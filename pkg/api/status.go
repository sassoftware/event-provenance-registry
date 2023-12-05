// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/sassoftware/event-provenance-registry/pkg/config"
	"github.com/sassoftware/event-provenance-registry/pkg/status"
)

// CheckLiveness() is a method of the `Server` struct. It
// returns an `http.HandlerFunc` function that handles the liveness check endpoint.
func (s *Server) CheckLiveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, `{"alive"":true}`)
	}
}

// CheckReadiness() is a method of the `Server` struct. It
// returns an `http.HandlerFunc` function that handles the readiness check endpoint. When this endpoint
// is called, it will respond with a JSON object `{"ready": true}` using the `render.JSON` function.
func (s *Server) CheckReadiness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, `{"ready":true}`)
	}
}

// CheckStatus() is a method of the `Server` struct. It returns
// an `http.HandlerFunc` function that handles the status check endpoint. When this endpoint is called,
// it will respond with a JSON object `{"ready": true}` using the `render.JSON` function.
func (s *Server) CheckStatus(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := status.New(cfg)
		render.JSON(w, r, s)
	}
}
