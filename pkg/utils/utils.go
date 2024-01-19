// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

// NowRFC3339 returns an RFC3339 format string
func NowRFC3339() string {
	return time.Now().Format(time.RFC3339)
}

// GetEnv returns an env variable value or a default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// NewULID returns a ULID.
func NewULID() (ulid.ULID, error) {
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	if err != nil {
		return id, fmt.Errorf("NewULID Failed: %s", err)
	}
	return id, err
}

// NewULIDAsString returns a ULID string.
func NewULIDAsString() string {
	id, _ := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	return id.String()
}

// NewUUID returns a UUID.
func NewUUID() (uuid.UUID, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return id, fmt.Errorf("NewUUID Failed: %s", err)
	}
	return id, err
}

// NewUUIDAsString returns a UUID as a String.
func NewUUIDAsString() string {
	id, _ := uuid.NewUUID()
	return id.String()
}

// BCrypt is a new type to make handling of the api key easier
type BCrypt struct {
	key []byte
}

// NewBCrypt is used for storing the key as bytes
func NewBCrypt(key string) BCrypt {
	return BCrypt{[]byte(key)}
}

// Hash uses GenerateFromPassword to hash with a salt the key
func (bc BCrypt) Hash() (string, error) {
	// NOTE: GenerateFromPassword generates a salt automatically
	hash, err := bcrypt.GenerateFromPassword(bc.key, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Equal uses CompareHashAndPassword to valid a key and Hash compare
func (bc BCrypt) Equal(hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash), bc.key)
	if err == nil {
		return true, nil
	}
	return false, err
}

// IntInSlice is used to find a int in a list
func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// StringInSlice is used to find a string in a list
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// DeepEqualStringArray compares []string to []string and returns bool
func DeepEqualStringArray(first []string, second []string) bool {
	if len(first) != len(second) {
		return false
	}
	for _, v := range first {
		if !StringInSlice(v, second) {
			return false
		}
	}
	return true
}

// Seed struct for computing the fingerprint
type Seed struct {
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

// Fingerprint creates a fingerprint for an Event Receiver and a Receiver Group
// Order matters
// ["type", "description", "name", "version"]
func (g *Seed) Fingerprint() string {
	sep := " "
	seed := "v1" + sep + g.Type + sep + g.Description + sep + g.Name + sep + g.Version
	sum := sha256.Sum256([]byte(seed))
	return fmt.Sprintf("%x", sum)
}

// MustGetLogger for logging
func MustGetLogger(name, module string) *logr.Logger {
	return getLogger(name, module)
}

func getLogger(name, module string) *logr.Logger {
	if name == "" {
		name = "logger"
	}
	if module == "" {
		module = "utils.logger"
	}

	logLevel, err := strconv.Atoi(GetEnv("LOG_LEVEL", "1"))
	if err != nil {
		logLevel = int(zerolog.InfoLevel) // default to INFO
	}
	zerologr.SetMaxV(1)

	zl := zerolog.New(os.Stderr).Level(zerolog.Level(logLevel)).With().Timestamp().Logger()

	logger := zerologr.New(&zl).WithName(name).WithValues("module", module)

	return &logger
}
