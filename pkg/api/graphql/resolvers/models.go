// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package resolvers

import (
	"github.com/graph-gophers/graphql-go"
)

type FindEventInput struct {
	ID              *graphql.ID
	Name            graphql.NullString
	Version         graphql.NullString
	Release         graphql.NullString
	PlatformID      graphql.NullString
	Package         graphql.NullString
	Success         graphql.NullBool
	EventReceiverID *graphql.ID
}

func (f FindEventInput) toMap() map[string]any {
	m := map[string]any{}
	if f.ID != nil {
		m["id"] = *f.ID
	}
	if f.Name.Set {
		m["name"] = f.Name.Value
	}
	if f.Version.Set {
		m["version"] = f.Version.Value
	}
	if f.Release.Set {
		m["release"] = f.Release.Value
	}
	if f.PlatformID.Set {
		m["platform_id"] = f.PlatformID.Value
	}
	if f.Package.Set {
		m["package"] = f.Package.Value
	}
	if f.Success.Set {
		m["success"] = f.Success.Value
	}
	if f.EventReceiverID != nil {
		m["event_receiver_id"] = *f.EventReceiverID
	}
	return m
}

type FindEventReceiverInput struct {
	ID      *graphql.ID
	Name    graphql.NullString
	Type    graphql.NullString
	Version graphql.NullString
}

func (f FindEventReceiverInput) toMap() map[string]any {
	m := map[string]any{}
	if f.ID != nil {
		m["id"] = *f.ID
	}
	if f.Name.Set {
		m["name"] = f.Name.Value
	}
	if f.Type.Set {
		m["type"] = f.Type.Value
	}
	if f.Version.Set {
		m["version"] = f.Version.Value
	}
	return m
}

type FindEventReceiverGroupInput struct {
	ID      *graphql.ID
	Name    graphql.NullString
	Type    graphql.NullString
	Version graphql.NullString
}

func (f FindEventReceiverGroupInput) toMap() map[string]any {
	m := map[string]any{}
	if f.ID != nil {
		m["id"] = *f.ID
	}
	if f.Name.Set {
		m["name"] = f.Name.Value
	}
	if f.Type.Set {
		m["type"] = f.Type.Value
	}
	if f.Version.Set {
		m["version"] = f.Version.Value
	}
	return m
}
