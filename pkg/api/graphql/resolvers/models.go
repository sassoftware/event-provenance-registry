// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package resolvers

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
)

type EventInput struct {
	Name            string
	Version         string
	Release         string
	PlatformID      string
	Package         string
	Description     string
	Payload         types.JSON
	Success         bool
	EventReceiverID graphql.ID
}

type EventReceiverInput struct {
	Name        string
	Type        string
	Version     string
	Description string
	Schema      types.JSON
}

type EventReceiverGroupInput struct {
	Name             string
	Type             string
	Version          string
	Description      string
	EventReceiverIDs []graphql.ID
}

type FindEventInput struct {
	ID              *graphql.ID
	Name            *string
	Version         *string
	Release         *string
	PlatformID      *string
	Package         *string
	Success         *bool
	EventReceiverID *graphql.ID
}

type FindEventReceiverInput struct {
	ID      *graphql.ID
	Name    *string
	Type    *string
	Version *string
}

type FindEventReceiverGroupInput struct {
	ID      *graphql.ID
	Name    *string
	Type    *string
	Version *string
}
