//go:generate go run github.com/99designs/gqlgen generate

package graph

import "gitlab.sas.com/async-event-infrastructure/server/pkg/storage"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Database *storage.Database
}
