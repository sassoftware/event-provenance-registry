// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"fmt"

	"github.com/sassoftware/event-provenance-registry/pkg/config"
)

var (
	// AppName is application name
	AppName = "server"
	// AppVersion is the application version
	AppVersion = config.Version
	// AppRelease is the application release
	AppRelease = config.Commit
)

// Version struct for storing Version info
type Version struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	Release string `json:"release" yaml:"release"`
}

// GetVersion returns name-version-release as a string
func GetVersion() string {
	return AppName + "-" + AppVersion + "-" + AppRelease
}

// GetVersionJSON returns version info as JSON
func GetVersionJSON() string {
	return fmt.Sprintf(`{"name": "%s", "version": "%s", "release": "%s"}`, AppName, AppVersion, AppRelease)
}

// NewVersion returns a populated Version struct
func NewVersion() *Version {
	return &Version{
		Name:    AppName,
		Version: AppVersion,
		Release: AppRelease,
	}
}
