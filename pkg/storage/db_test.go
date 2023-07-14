// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"log"
	"testing"
)

func Test_InitDB(t *testing.T) {
	db := connToDB()

	err := db.SyncSchema()
	if err != nil {
		log.Fatal(err)
	}
}

func Test_CreateEvent(t *testing.T) {
	_ = connToDB()

	// ...
}

func Test_CreateEventReceiver(t *testing.T) {
	db := connToDB()

	eventReceiver := EventReceiver{
		ID:          "1",
		Name:        "name",
		Type:        "type",
		Version:     "version",
		Description: "description",
		Enabled:     true,
	}

	_, err := CreateEventReceiver(db.Client, eventReceiver)
	if err != nil {
		log.Fatal(err)
	}
}

func Test_CreateEventReceiverGroup(t *testing.T) {
	_ = connToDB()

	// ...
}

func connToDB() *Database {
	db, err := New("localhost", "postgres", "", "", "postgres")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
