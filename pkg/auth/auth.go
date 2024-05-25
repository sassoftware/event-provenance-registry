// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"log/slog"
	"net/http"
	"strings"
)

type Middleware struct {
	verifier *oidc.IDTokenVerifier
}

func NewMiddleware(ctx context.Context) (*Middleware, error) {
	provider, err := oidc.NewProvider(ctx, "http://localhost:8083/realms/test")
	if err != nil {
		return nil, err
	}

	clientID := "epr-client-id"
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	verifier := provider.Verifier(oidcConfig)

	return &Middleware{
		verifier: verifier,
	}, nil
}

func (m *Middleware) Handle() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			bearer := BearerToken(r)

			if bearer == "" {
				slog.Error("bearer token missing")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			//	authHeader := r.Header.Get("authorization")
			//	split := strings.Split(authHeader, " ")
			//	if len(split) != 2 {
			//		slog.Error("more authorization header pieces than expected")
			//		w.WriteHeader(http.StatusBadRequest)
			//		return
			//	}
			//	if !strings.EqualFold(split[0], "bearer") {
			//		slog.Error("unexpected authorization header type", "type", split[0])
			//	}
			//
			//	idToken, err := verifier.Verify(r.Context(), split[1])
			//	if err != nil {
			//		slog.Error("authorization failed", "error", err)
			//		w.WriteHeader(http.StatusUnauthorized)
			//		return
			//	}

			_, err := m.verifier.Verify(r.Context(), bearer)
			if err != nil {
				slog.Error(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// TODO: validate any roles and stuff.

			// TODO: do we need this context key?
			//ctx := context.WithValue(r.Context(), "", token)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func BearerToken(req *http.Request) string {
	auth := req.Header.Get("Authorization")
	const prefix = "Bearer "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return ""
	}
	return auth[len(prefix):]
}
