// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/jackc/pgconn"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

var logger = utils.MustGetLogger("db", "pkg.storage")

type Database struct {
	Client *gorm.DB
}

func New(host, user, pass, sslMode, database string, port int) (*Database, error) {
	glog := gormlog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlog.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlog.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	if sslMode != "" {
		sslMode = fmt.Sprintf("sslmode=%s", sslMode)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d %v TimeZone=EST", host, user, pass, database, port, sslMode)
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: glog})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database. err %s", err)
	}
	return &Database{Client: client}, err
}

func (db *Database) SyncSchema() error {
	return db.Client.AutoMigrate(
		new(Event),
		new(EventReceiver),
		new(EventReceiverGroup),
		new(EventReceiverGroupToEventReceiver),
	)
}

func CreateEvent(tx *gorm.DB, event Event) (*Event, error) {
	event.ID = graphql.ID(utils.NewULIDAsString())

	if !receiversExist(tx, []string{string(event.EventReceiverID)}) {
		return nil, fmt.Errorf("cannot create event for non-existent receiver '%s'", event.EventReceiverID)
	}

	result := tx.Create(&event)
	if result.Error != nil {
		return nil, pgError(result.Error)
	}
	return &event, nil
}

func FindEvent(tx *gorm.DB, id graphql.ID) (*Event, error) {
	var event Event
	result := tx.Model(&Event{}).First(&event, &Event{ID: id})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("event %s not found", id)
		}
		return nil, pgError(result.Error)
	}
	return &event, nil
}

func CreateEventReceiver(tx *gorm.DB, eventReceiver EventReceiver) (*EventReceiver, error) {
	eventReceiver.ID = graphql.ID(utils.NewULIDAsString())

	// TODO: set fingerprint

	result := tx.Create(&eventReceiver)
	if result.Error != nil {
		return nil, pgError(result.Error)
	}
	return &eventReceiver, nil
}

func FindEventReceiver(tx *gorm.DB, id graphql.ID) (*EventReceiver, error) {
	var eventReciever EventReceiver
	result := tx.Model(&EventReceiver{}).First(&eventReciever, &EventReceiver{ID: id})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("eventReciever %s not found", id)
		}
		return nil, pgError(result.Error)
	}
	return &eventReciever, nil
}

func CreateEventReceiverGroup(tx *gorm.DB, receiverIDs []string, group EventReceiverGroup) (*EventReceiverGroup, error) {
	group.ID = graphql.ID(utils.NewULIDAsString())

	if !receiversExist(tx, receiverIDs) {
		return nil, fmt.Errorf("not all receivers exist")
	}

	// TODO: find a way to make these two creations atomic
	result := tx.Create(&group)
	if result.Error != nil {
		return nil, pgError(result.Error)
	}

	// TODO: This is a horrible way to do this. Figure out something better.
	for _, id := range receiverIDs {
		recToGroup := EventReceiverGroupToEventReceiver{
			EventReceiverGroupID: group.ID,
			EventReceiverID:      graphql.ID(id),
		}
		tx.Create(&recToGroup)
		if result.Error != nil {
			return nil, pgError(result.Error)
		}
	}

	return &group, nil
}

func FindEventReceiverGroup(tx *gorm.DB, id graphql.ID) (*EventReceiverGroup, error) {
	var eventRecieverGroup EventReceiverGroup
	result := tx.Model(&EventReceiverGroup{}).First(&eventRecieverGroup, &EventReceiverGroup{ID: id})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("eventRecieverGroup %s not found", id)
		}
		return nil, pgError(result.Error)
	}
	return &eventRecieverGroup, nil
}

func updateRecord(tx *gorm.DB, record any) error {
	result := tx.Save(record)
	if result.Error != nil {
		return pgError(result.Error)
	}
	return nil
}

func deleteRecord(tx *gorm.DB, record any) error {
	result := tx.Delete(record)
	if result.Error != nil {
		return pgError(result.Error)
	}
	return nil
}

func pgError(err error) error {
	switch err := err.(type) {
	case *pgconn.PgError:
		return fmt.Errorf("err: %s. detail: %s. code: %s", err.Message, err.Detail, err.SQLState())
	default:
		return err
	}
}

func receiversExist(tx *gorm.DB, ids []string) bool {
	var recs []EventReceiver
	for _, id := range ids {
		rec := EventReceiver{
			ID: graphql.ID(id),
		}
		recs = append(recs, rec)
	}
	result := tx.Find(&recs)
	if result.Error != nil {
		logger.V(1).Error(result.Error, "unable to obtain receiver")
		return false
	}

	if result.RowsAffected != int64(len(ids)) {
		//TODO: return data on missing receivers
		logger.V(1).Info("some event receivers do not exist")
		return false
	}
	return true
}
