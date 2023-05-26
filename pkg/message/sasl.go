// Copyright © 2021, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package message

import (
	"os"
)

// Mechanism represents the SASL Authentication mechanism being used
type Mechanism int

const (
	NONE Mechanism = iota
	PLAIN
	SCRAM
	OAUTH2 // OAUTH2 is not currently supported, but will be
)

type SASLAuthentication struct {
	Mechanism Mechanism
	Username  string
	Password  string
}

// SASLEnabled returns true if SASL_USERNAME, SASL_PASSWORD, and a SASL_MECHANISM are set
func (s *SASLAuthentication) SASLEnabled() bool {
	return s.Username != "" && s.Password != "" && s.Mechanism != NONE
}

func getMechanism() Mechanism {
	switch os.Getenv("SASL_MECHANISM") {
	case "PLAIN":
		return PLAIN
	case "SCRAM":
		return SCRAM
	case "OAUTH2":
		return OAUTH2
	default:
		return NONE
	}
}

// GetSASLCredentials returns SASLCredentials struct
func GetSASLAuthentication() *SASLAuthentication {
	return &SASLAuthentication{
		Username:  os.Getenv("SASL_USERNAME"),
		Password:  os.Getenv("SASL_PASSWORD"),
		Mechanism: getMechanism(),
	}
}
