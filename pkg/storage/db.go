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
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
	"github.com/sassoftware/event-provenance-registry/pkg/utils"
	"github.com/xeipuuv/gojsonschema"
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

// CreateEvent creates and event record in the database. Throws an error if the event receiver does not exist or if the
// event payload does not match the receiver schema.
func CreateEvent(tx *gorm.DB, event Event) (*Event, error) {
	receivers, err := FindEventReceiver(tx, event.EventReceiverID)
	if err != nil {
		return nil, fmt.Errorf("could not validate event schema. %w", err)
	}

	if len(*receivers) > 1 {
		return nil, fmt.Errorf("more than one receiver was found with ID %s", event.EventReceiverID)
	}
	receiver := (*receivers)[0]

	if err := validateReceiverSchema(receiver.Schema.String(), event.Payload); err != nil {
		return nil, err
	}
	event.ID = graphql.ID(utils.NewULIDAsString())

	results := tx.Create(&event)
	if results.Error != nil {
		return nil, pgError(results.Error)
	}
	event.EventReceiver = *receiver
	return &event, nil
}

func FindEvent(tx *gorm.DB, id graphql.ID) (*[]*Event, error) {
	var events []*Event
	result := tx.Model(&Event{}).Preload("EventReceiver").First(&events, &Event{ID: id})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("event %s not found", id)
		}
		return nil, pgError(result.Error)
	}
	return &events, nil
}

func CreateEventReceiver(tx *gorm.DB, eventReceiver EventReceiver) (*EventReceiver, error) {
	eventReceiver.ID = graphql.ID(utils.NewULIDAsString())

	seed := utils.Seed{
		Name:        eventReceiver.Name,
		Type:        eventReceiver.Type,
		Version:     eventReceiver.Version,
		Description: eventReceiver.Description,
	}
	eventReceiver.Fingerprint = seed.Fingerprint()

	result := tx.Create(&eventReceiver)
	if result.Error != nil {
		return nil, pgError(result.Error)
	}
	return &eventReceiver, nil
}

// FindEventReceiver tries to find an event receiver by ID.
func FindEventReceiver(tx *gorm.DB, id graphql.ID) (*[]*EventReceiver, error) {
	var eventReceivers []*EventReceiver
	result := tx.Model(&EventReceiver{}).First(&eventReceivers, &EventReceiver{ID: id})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("eventReceiver %s not found", id)
		}
		return nil, pgError(result.Error)
	}
	return &eventReceivers, nil
}

func CreateEventReceiverGroup(tx *gorm.DB, eventReceiverGroup EventReceiverGroup) (*EventReceiverGroup, error) {
	eventReceiverGroup.ID = graphql.ID(utils.NewULIDAsString())

	// create our EventReceiverGroupToEventReceivers
	eventReceiverGroupToEventReceivers := []*EventReceiverGroupToEventReceiver{}
	for _, eventReceiverID := range eventReceiverGroup.EventReceiverIDs {
		eventReceiverGroupToEventReceivers = append(eventReceiverGroupToEventReceivers, &EventReceiverGroupToEventReceiver{
			EventReceiverID:      eventReceiverID,
			EventReceiverGroupID: eventReceiverGroup.ID,
		})
	}

	// create both EventReceiverGroup and EventReceiverGroupToEventReceiver in a single transaction
	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&eventReceiverGroup)
		if result.Error != nil {
			return pgError(result.Error)
		}

		result = tx.CreateInBatches(eventReceiverGroupToEventReceivers, len(eventReceiverGroupToEventReceivers))
		if result.Error != nil {
			return pgError(result.Error)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &eventReceiverGroup, nil
}
func FindEventReceiverGroup(tx *gorm.DB, id graphql.ID) (*[]*EventReceiverGroup, error) {
	var eventReceiverGroup EventReceiverGroup
	result := tx.Model(&EventReceiverGroup{}).First(&eventReceiverGroup, &EventReceiverGroup{ID: id})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("eventReceiverGroup %s not found", id)
		}
		return nil, pgError(result.Error)
	}

	result = tx.Model(&EventReceiverGroupToEventReceiver{}).
		Select("event_receiver_id").
		Find(&eventReceiverGroup.EventReceiverIDs, &EventReceiverGroupToEventReceiver{EventReceiverGroupID: id})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("eventReceiverGroup %s not found in EventReceiverGroupToEventReceiver", id)
		}
		return nil, pgError(result.Error)
	}

	return &[]*EventReceiverGroup{&eventReceiverGroup}, nil
}

func SetEventReceiverGroupEnabled(tx *gorm.DB, id graphql.ID, enabled bool) error {
	result := tx.Model(&EventReceiverGroup{ID: id}).Update("enabled", enabled)
	return pgError(result.Error)
}

func pgError(err error) error {
	switch err := err.(type) {
	case *pgconn.PgError:
		return fmt.Errorf("err: %s. detail: %s. code: %s", err.Message, err.Detail, err.SQLState())
	default:
		return err
	}
}

func validateReceiverSchema(schema string, eventPayload types.JSON) error {
	loader := gojsonschema.NewStringLoader(schema)
	sch, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return err
	}
	jsResult, err := sch.Validate(gojsonschema.NewStringLoader(eventPayload.String()))
	if err != nil {
		return err
	}
	if !jsResult.Valid() {
		msg := "event payload did not match event receiver schema"
		err := errors.New(msg)
		for _, e := range jsResult.Errors() {
			err = errors.Join(err, errors.New(e.String()))
		}
		logger.Error(err, "invalid schema")
		return err
	}

	return nil
}
