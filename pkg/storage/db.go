// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/jackc/pgconn"
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
	eprErrors "github.com/sassoftware/event-provenance-registry/pkg/errors"
	"github.com/sassoftware/event-provenance-registry/pkg/utils"
	"github.com/xeipuuv/gojsonschema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

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
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: glog, PrepareStmt: true})
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
	receivers, err := FindEventReceiverByID(tx, event.EventReceiverID)
	if err != nil {
		switch err.(type) {
		case eprErrors.MissingObjectError:
			return nil, eprErrors.InvalidInputError{Msg: "receiver for event does not exist"}
		default:
			return nil, err
		}
	}

	if len(receivers) > 1 {
		return nil, fmt.Errorf("more than one receiver was found with ID %s", event.EventReceiverID)
	}
	receiver := receivers[0]

	if err := validateReceiverSchema(receiver.Schema.String(), event.Payload); err != nil {
		return nil, eprErrors.InvalidInputError{Msg: err.Error()}
	}
	event.ID = graphql.ID(utils.NewULIDAsString())

	results := tx.Create(&event)
	if results.Error != nil {
		return nil, pgError(results.Error)
	}
	event.EventReceiver = receiver
	return &event, nil
}

func FindEventByID(tx *gorm.DB, id graphql.ID) ([]Event, error) {
	return FindEvent(tx, map[string]any{"id": id})
}

func FindEvent(tx *gorm.DB, e map[string]any) ([]Event, error) {
	var events []Event
	result := tx.Model(&Event{}).Preload("EventReceiver").Where(e).Find(&events)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, eprErrors.MissingObjectError{Msg: fmt.Sprintf("event %+v not found", e)}
		}
		return nil, pgError(result.Error)
	}
	return events, nil
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
func FindEventReceiverByID(tx *gorm.DB, id graphql.ID) ([]EventReceiver, error) {
	return FindEventReceiver(tx, map[string]any{"id": id})
}

// FindEventReceiver tries to find an event receiver by ID.
func FindEventReceiver(tx *gorm.DB, er map[string]any) ([]EventReceiver, error) {
	var eventReceivers []EventReceiver
	result := tx.Model(&EventReceiver{}).Where(er).Find(&eventReceivers)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, eprErrors.MissingObjectError{Msg: fmt.Sprintf("eventReceiver %+v not found", er)}
		}
		return nil, pgError(result.Error)
	}
	return eventReceivers, nil
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

func FindEventReceiverGroupByID(tx *gorm.DB, id graphql.ID) ([]EventReceiverGroup, error) {
	return FindEventReceiverGroup(tx, map[string]any{"id": id})
}

func FindEventReceiverGroup(tx *gorm.DB, erg map[string]any) ([]EventReceiverGroup, error) {
	var eventReceiverGroups []EventReceiverGroup
	result := tx.Model(&EventReceiverGroup{}).Where(erg).Find(&eventReceiverGroups)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, eprErrors.MissingObjectError{Msg: fmt.Sprintf("eventReceiverGroup %+v not found", erg)}
		}
		return nil, pgError(result.Error)
	}

	for i := range eventReceiverGroups {
		// need indirection so db query can modify array contents
		eventReceiverGroup := &eventReceiverGroups[i]
		result = tx.Model(&EventReceiverGroupToEventReceiver{}).
			Select("event_receiver_id").
			Find(&eventReceiverGroup.EventReceiverIDs, &EventReceiverGroupToEventReceiver{EventReceiverGroupID: eventReceiverGroup.ID})
		if result.Error != nil {
			return nil, pgError(result.Error)
		}
	}

	return eventReceiverGroups, nil
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
		slog.Error("invalid schema", "error", err)
		return err
	}

	return nil
}

//go:embed megaQuery.sql
var megaQuery string

// Data represents the EventReceiverGroup data that comes back from
// the mega query. It is necessary to allow us to automatically
// insert event_receiver_ids by overriding the json tag
type TriggeredEventReceiverGroups struct {
	EventReceiverGroup

	// Database returns json array instead of a string
	EventReceiverIDs types.JSON `json:"event_receiver_ids"`
}

func FindTriggeredEventReceiverGroups(tx *gorm.DB, event Event) ([]EventReceiverGroup, error) {
	var triggeredEventReceiverGroups []TriggeredEventReceiverGroups
	result := tx.Model("EventReceiverGroup").Raw(megaQuery,
		sql.Named("name", event.Name),
		sql.Named("version", event.Version),
		sql.Named("release", event.Release),
		sql.Named("platform_id", event.PlatformID),
		sql.Named("package", event.Package),
		sql.Named("event_receiver_id", event.EventReceiverID)).
		Scan(&triggeredEventReceiverGroups)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("eventReceiverGroup %s not found", event.EventReceiverID)
		}
		return nil, pgError(result.Error)
	}
	if result.RowsAffected < 1 {
		return nil, nil
	}

	eventReceiverGroups := []EventReceiverGroup{}
	for _, triggeredEventReceiverGroup := range triggeredEventReceiverGroups {
		var eventReceiverIDs []graphql.ID
		err := json.Unmarshal([]byte(triggeredEventReceiverGroup.EventReceiverIDs.JSON), &eventReceiverIDs)
		if err != nil {
			return nil, err
		}
		eventReceiverGroup := EventReceiverGroup{
			ID:               triggeredEventReceiverGroup.ID,
			Name:             triggeredEventReceiverGroup.Name,
			Version:          triggeredEventReceiverGroup.Version,
			Description:      triggeredEventReceiverGroup.Description,
			Enabled:          triggeredEventReceiverGroup.Enabled,
			EventReceiverIDs: eventReceiverIDs,
			CreatedAt:        triggeredEventReceiverGroup.CreatedAt,
			UpdatedAt:        triggeredEventReceiverGroup.UpdatedAt,
		}
		eventReceiverGroups = append(eventReceiverGroups, eventReceiverGroup)
	}
	return eventReceiverGroups, nil
}
