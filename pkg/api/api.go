// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/config"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/server"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/status"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/utils"
)

var logger = utils.MustGetLogger("server", "server.api")

// InitializeAPI starts the database, kafka message producer, middleware, and endpoints
func InitializeAPI(_ context.Context, cfg *config.Config) (*chi.Mux, *storage.Database, error) {
	// Create a new router
	router := chi.NewRouter()

	// Create a new connection to our pg database
	db, err := storage.New(cfg.DB.Host, cfg.DB.User, cfg.DB.Pass, cfg.DB.SSLMode, cfg.DB.Name)
	if err != nil {
		log.Fatal(err)
	}

	// db.SetMaxOpenConns(cfg.DB.MaxConnections)
	// db.SetMaxIdleConns(cfg.DB.IdleConnections)
	// db.SetConnMaxLifetime(time.Duration(cfg.DB.ConnectionLife) * time.Minute)

	//  TODO: add this stuff.
	// schema, err := graph.NewSchema(db)
	// if err != nil {
	// 	logger.Error(err, "error creating schema")
	// 	return nil, nil, err
	// }
	//
	// schema.AddExtensions(graph.NewLastPageExt())
	// schema.AddExtensions(graph.NewTotalExt())
	//
	// //  Create a server struct that holds a pointer to our database as well
	// //  as the address of our graphql schema
	// authorizer, err := auth.NewAuthorizer(ctx, cfg.Auth)
	// if err != nil {
	// 	return nil, nil, err
	// }
	s := server.New() // cfg, &schema, db, cfg.Kafka.MsgChannel

	//  Add some middleware to our router
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
		Debug:            true,
	})

	router.Route("/api", func(r chi.Router) {
		r.Get("/", s.ServeOpenAPIDoc(cfg.ResourceDir))
		r.Route("/v1", func(r chi.Router) {
			r.Use(crs.Handler)
			if cfg.Verbose {
				httpLogger := httplog.NewLogger("server-http-logger", httplog.Options{
					JSON: true,
					Tags: map[string]string{
						"version": status.AppVersion,
						"release": status.AppRelease,
					},
				})
				r.Use(httplog.RequestLogger(httpLogger))
			}
			r.Get("/openapi", s.ServeOpenAPIDoc(cfg.ResourceDir))
			// REST endpoints
			r.Route("/events", func(r chi.Router) {
				r.With(s.Paginate).With(s.Sorting).Get("/", s.GetEvents())
				r.With(s.Paginate).With(s.Sorting).Head("/", s.GetEvents())
				r.Post("/", s.CreateEvent())
				r.Route("/{eventID}", func(r chi.Router) {
					r.Get("/", s.GetEventByID())
				})
			})
			r.Route("/receivers", func(r chi.Router) {
				r.Post("/", s.CreateReceiver())
				r.With(s.Paginate).With(s.Sorting).Get("/", s.GetReceivers())
				r.With(s.Paginate).With(s.Sorting).Head("/", s.GetReceivers())
				r.Route("/{receiverID}", func(r chi.Router) {
					r.Get("/", s.GetReceiverByID())
				})
			})
			r.Route("/groups", func(r chi.Router) {
				r.Post("/", s.CreateGroup())
				r.With(s.Paginate).With(s.Sorting).Get("/", s.GetGroups())
				r.With(s.Paginate).With(s.Sorting).Head("/", s.GetGroups())
				r.Route("/{groupID}", func(r chi.Router) {
					r.Get("/", s.GetGroupByID())
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
	if cfg.Debug {
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

		r.Post("/", s.GraphQLHandler())
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

	FileServer(router, "/resources", cfg.ResourceDir)
	return router, db, nil
}

// FileServer is serving static files
func FileServer(r chi.Router, public string, static string) {
	if strings.ContainsAny(public, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	root, _ := filepath.Abs(static)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		panic("Static Documents Directory Not Found")
	}

	fs := http.StripPrefix(public, http.FileServer(http.Dir(root)))

	if public != "/" && public[len(public)-1] != '/' {
		r.Get(public, http.RedirectHandler(public+"/", http.StatusMovedPermanently).ServeHTTP)
		public += "/"
	}

	// authorizer.RequireRole(multitoken.RoleUnauthenticated) https://rndjira.sas.com/browse/CPIPE-89
	r.With().Get(public+"*", func(w http.ResponseWriter, r *http.Request) {
		file := strings.Replace(r.RequestURI, public, "/", 1)
		if _, err := os.Stat(root + file); os.IsNotExist(err) {
			http.ServeFile(w, r, path.Join(root, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})
}
