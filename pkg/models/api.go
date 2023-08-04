// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package models

import "gitlab.sas.com/async-event-infrastructure/server/pkg/storage"

// GroupInput rest representation of the data for a storage.EventReceiverGroup
type GroupInput struct {
	storage.EventReceiverGroup
	EventReceiverIDs []string `json:"event_receiver_ids"`
}
