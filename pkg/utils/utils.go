// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"golang.org/x/crypto/bcrypt"
)

// GetEnv returns an env variable value or a default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetEnvsByPrefix finds all ENV vars that start with prefix
// GetEnvsByPrefix func takes no as input and returns prefix string, strip bool map[string]string
func GetEnvsByPrefix(prefix string, strip bool) map[string]string {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if strings.HasPrefix(pair[0], prefix) {
			if len(pair[1]) > 0 {
				k := pair[0]
				if strip {
					k = strings.Split(pair[0], prefix+"_")[1]
				}
				envs[k] = pair[1]
			}
		}
	}
	return envs
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
	Name        string   `json:"name,omitempty"`
	Action      string   `json:"action,omitempty"`
	Version     string   `json:"version,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Fingerprint creates a fingerprint for a Gate or Stage
// Order matters
// ["action", "description", "name", "tags", "version"]
func (g *Seed) Fingerprint() string {
	sep := " "
	seed := "v1" + sep + g.Action + sep + g.Description + sep + g.Name
	for _, x := range g.Tags {
		seed = seed + sep + x
	}
	seed = seed + sep + g.Version
	sum := sha256.Sum256([]byte(seed))
	return fmt.Sprintf("%x", sum)
}
