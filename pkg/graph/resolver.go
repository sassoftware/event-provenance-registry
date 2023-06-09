//go:generate go run github.com/99designs/gqlgen generate

package graph

import "gitlab.sas.com/async-event-infrastructure/server/pkg/db"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	Database *db.Database 
}
