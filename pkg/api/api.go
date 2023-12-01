// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sassoftware/event-provenance-registry/pkg/config"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/status"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/sassoftware/event-provenance-registry/pkg/utils"
)

var logger = utils.MustGetLogger("server", "server.api")

// Initialize starts the database, kafka message producer, middleware, and endpoints
func Initialize(ctx context.Context, cfg *config.Config) (*chi.Mux, error) {
	if cfg == nil {
		return nil, fmt.Errorf("no config provided")
	}

	//// Create a new connection to our pg database
	// db, err := storage.New(cfg.DB.Host, cfg.DB.User, cfg.DB.Pass, cfg.DB.SSLMode, cfg.DB.Name)
	// if err != nil {
	//	log.Fatal(err)
	// }

	// Create a new connection to our pg database
	connection, err := storage.New(cfg.Storage.Host, cfg.Storage.User, cfg.Storage.Pass, cfg.Storage.SSLMode, cfg.Storage.Name, cfg.Storage.Port)
	if err != nil {
		log.Fatal(err)
	}

	err = connection.SyncSchema()
	if err != nil {
		log.Fatal(err)
	}

	// set up kafka
	kafkaCfg, err := message.NewConfig(cfg.Kafka.Version)
	if err != nil {
		log.Fatal(err)
	}

	cfg.Kafka.Producer, err = message.NewProducer(cfg.Kafka.Peers, kafkaCfg)
	if err != nil {
		log.Fatal(err)
	}

	cfg.Kafka.Producer.ConsumeSuccesses()
	cfg.Kafka.Producer.ConsumeErrors()

	s, err := New(ctx, connection, cfg.Kafka)
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
		r.Get("/", s.Rest.ServeOpenAPIDoc(cfg.Server.ResourceDir))
		r.Route("/v1", func(r chi.Router) {
			r.Use(crs.Handler)
			if cfg.Server.VerboseAPI {
				httpLogger := httplog.NewLogger("server-http-logger", httplog.Options{
					JSON: true,
					Tags: map[string]string{
						"version": status.AppVersion,
						"release": status.AppRelease,
					},
				})
				r.Use(httplog.RequestLogger(httpLogger))
			}
			r.Get("/openapi", s.Rest.ServeOpenAPIDoc(cfg.Server.ResourceDir))
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
		r.Get("/status", s.CheckStatus())
	})

	// turn on the profiler in debug mode
	if cfg.Server.Debug {
		// profiler
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
		r.Post("/query", s.GraphQL.GraphQLHandler(connection, cfg.Kafka))
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
