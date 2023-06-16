// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Client *gorm.DB
}

func NewDB(host, user, pass, sslMode, database string) (*Database, error) {
	glog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	if sslMode != "" {
		sslMode = fmt.Sprintf("sslmode=%s", sslMode)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 %s TimeZone=EST", host, user, pass, database, sslMode)
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
	result := tx.Create(&event)
	if result.Error != nil {
		return nil, pgError(result.Error)
	}
	return &event, nil
}

func CreateEventReceiver(tx *gorm.DB, eventReceiver EventReceiver) (*EventReceiver, error) {
	result := tx.Create(&eventReceiver)
	if result.Error != nil {
		return nil, pgError(result.Error)
	}
	return &eventReceiver, nil
}

func CreateEventReceiverGroup(tx *gorm.DB, eventReceiverGroup EventReceiverGroup) (*EventReceiverGroup, error) {
	result := tx.Create(&eventReceiverGroup)
	if result.Error != nil {
		return nil, pgError(result.Error)
	}
	return &eventReceiverGroup, nil
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
