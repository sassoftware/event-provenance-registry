// Copyright © 2021, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestVersion(t *testing.T) {
	ver := GetVersion()
	assert.Equal(t, ver, "server-dev-dirty")
	jsonvar := GetVersionJSON()
	assert.Equal(t, jsonvar, `{"name": "server", "version": "dev", "release": "dirty"}`)
	verStruct := NewVersion()
	assert.Equal(t, verStruct.Name, "server")
	assert.Equal(t, verStruct.Version, "dev")
	assert.Equal(t, verStruct.Release, "dirty")
}
