// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sassoftware/event-provenance-registry/pkg/config"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/status"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// Initialize starts the database, kafka message producer, middleware, and endpoints
func Initialize(db *storage.Database, msgProducer message.TopicProducer, cfg *config.ServerConfig) (*chi.Mux, error) {
	if cfg == nil {
		return nil, fmt.Errorf("no config provided")
	}

	s, err := New(db, msgProducer)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new router
	router := chi.NewRouter()
	// Add some middleware to our router
	router.Use(
		requestTimer(),
		requestCounter(),
		securityHeaders(),
		render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		middleware.Logger,       // log api request calls
		middleware.Compress(5),  // compress results, mostly gzipping assets and json
		middleware.StripSlashes, // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer,    // recover from panics without crashing server
		middleware.GetHead,      // route HEAD requests
		LogOrigin,
	)

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"X-Total-Count", "X-Last-Page"},
		MaxAge:           300,
		Debug:            false,
	})

	router.Route("/api", func(r chi.Router) {
		r.Get("/", s.Rest.ServeOpenAPIDoc(cfg.ResourceDir))
		r.Route("/v1", func(r chi.Router) {
			r.Use(crs.Handler)
			if cfg.VerboseAPI {
				httpLogger := httplog.NewLogger("server-http-logger", httplog.Options{
					JSON: true,
					Tags: map[string]string{
						"version": status.AppVersion,
						"release": status.AppRelease,
					},
				})
				r.Use(httplog.RequestLogger(httpLogger))
			}
			r.Get("/openapi", s.Rest.ServeOpenAPIDoc(cfg.ResourceDir))
			// REST endpoints
			r.Route("/events", func(r chi.Router) {
				r.Post("/", s.Rest.CreateEvent())
				r.Route("/{eventID}", func(r chi.Router) {
					r.Get("/", s.Rest.GetEventByID())
				})
			})
			r.Route("/receivers", func(r chi.Router) {
				r.Post("/", s.Rest.CreateReceiver())
				r.Route("/{receiverID}", func(r chi.Router) {
					r.Get("/", s.Rest.GetReceiverByID())
				})
			})
			r.Route("/groups", func(r chi.Router) {
				r.Post("/", s.Rest.CreateGroup())
				r.Put("/enable", s.Rest.SetGroupEnabled(true))
				r.Put("/disable", s.Rest.SetGroupEnabled(false))
				r.Route("/{groupID}", func(r chi.Router) {
					r.Get("/", s.Rest.GetGroupByID())
				})
			})
		})
	})

	router.Route("/healthz", func(r chi.Router) {
		r.Get("/liveness", s.CheckLiveness())
		r.Get("/readiness", s.CheckReadiness())
		r.Get("/status", s.CheckStatus(cfg))
	})

	// turn on the profiler in debug mode
	if cfg.Debug {
		// profiler
		slog.Debug("profiler enabled")
		router.Route("/", func(r chi.Router) {
			r.Mount("/debug", middleware.Profiler())
		})
	}

	// endpoint for serving Prometheus metrics
	router.Route("/metrics", func(r chi.Router) {
		r.Get("/", promhttp.Handler().(http.HandlerFunc))
	})

	// Separate, to ensure no authentication required.
	router.Route("/api/v1/graphql", func(r chi.Router) {
		r.Use(crs.Handler)
		r.Get("/", s.GraphQL.ServerGraphQLDoc())
		r.Post("/query", s.GraphQL.GraphQLHandler())
	})

	// Public Api Endpoints
	router.Group(func(r chi.Router) {
		crs := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET"},
			AllowedHeaders: []string{"Accept", "Content-Type"},
			ExposedHeaders: []string{"Link"},
			MaxAge:         300, // Maximum value not ignored by any of major browsers
		})
		// Add some middleware to our router
		r.Use(crs.Handler,
			render.SetContentType(render.ContentTypeHTML), // set content-type headers as text/html
			middleware.Logger,       // log api request calls
			middleware.Compress(5),  // compress results, mostly gzipping assets and json
			middleware.StripSlashes, // match paths with a trailing slash, strip it, and continue routing through the mux
			middleware.Recoverer,    // recover from panics without crashing server
		)
		// r.Get("/", s.GetIndexHTML())
		// r.Get("/status", s.GetServerStatusHTML())
	})
	return router, nil
}
