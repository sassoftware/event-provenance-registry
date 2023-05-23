// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os"
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("FOO", "1")
	foo := GetEnv("FOO", "2")
	assert.Assert(t, foo == "1")
	bar := GetEnv("BAR", "42")
	assert.Assert(t, bar == "42")
}

// TestGetEnvsByPrefix
func TestGetEnvsByPrefix(t *testing.T) {
	// 46 and 2
	os.Setenv("GET_ENV_PREFIX_FOO", "46")
	os.Setenv("GET_ENV_PREFIX_BAR", "2")
	prefix := "GET_ENV_PREFIX"
	tokens := GetEnvsByPrefix(prefix, true)
	assert.Assert(t, tokens["FOO"] == "46")
	assert.Assert(t, tokens["BAR"] == "2")
}

func TestNewULID(t *testing.T) {
	u, err := NewULID()
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, len(u.String()) == 26)
}

func TestNewULIDAsString(t *testing.T) {
	u := NewULIDAsString()
	assert.Assert(t, len(u) == 26)
}

func TestNewUUID(t *testing.T) {
	u, err := NewUUID()
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, len(u.String()) == 36)
}

func TestNewUUIDAsString(t *testing.T) {
	u := NewUUIDAsString()
	assert.Assert(t, len(u) == 36)
}

// TestNewBCryptWithApiKey validates the generation of a salted hash and verifies it against the original api key
func TestNewBCryptWithApiKey(t *testing.T) {
	bc := NewBCrypt("def456")
	hashedAPIKey, err := bc.Hash()
	assert.Assert(t, is.Nil(err))
	equal, err := bc.Equal(hashedAPIKey)
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, equal)
}

// TestNewBCryptWithDifferentApiKeys validates two different api key hashs aren't equal
func TestNewBCryptWithDifferentApiKeys(t *testing.T) {
	bc := NewBCrypt("def456")
	hashedAPIKey, err := bc.Hash()
	assert.Assert(t, is.Nil(err))
	_ = hashedAPIKey
	bc2 := NewBCrypt("jlp456")
	hashedAPIKey2, err := bc2.Hash()
	assert.Assert(t, is.Nil(err))
	equal, _ := bc.Equal(hashedAPIKey2)
	assert.Assert(t, !equal)
}

// TestIntInSlice validates IntInSlice works
func TestIntInSlice(t *testing.T) {
	list := []int{1, 42, 46, 2}
	result := IntInSlice(42, list)
	assert.Equal(t, result, true, "Error in IntInSlice")
	nr := IntInSlice(12, list)
	assert.Equal(t, nr, false, "Error in IntInSlice")
}

// TestStringInSlice validates StringInSlice works
func TestStringInSlice(t *testing.T) {
	list := []string{"foo", "bar", "caz"}
	result := StringInSlice("foo", list)
	assert.Equal(t, result, true, "Error in StringInSlice")
	nr := StringInSlice("Joe.Strummer", list)
	assert.Equal(t, nr, false, "Error in StringInSlice")
}

func TestDeepEqualStringArray(t *testing.T) {
	a := []string{"01DQDHG7GK3NQDREC3P47DJT48", "01DQWKGE76ABCYDZ07E6DZ2QC7"}
	b := []string{"01DQWKGE76ABCYDZ07E6DZ2QC7", "01DQDHG7GK3NQDREC3P47DJT48"}
	c := []string{"01DQWKGE76ABCYDZ07E6DZ2QC7", "foobar"}
	d := []string{"01DQWKGE76ABCYDZ07E6DZ2QC7"}
	assert.Assert(t, DeepEqualStringArray(a, b), "Error DeepEqualArray")
	assert.Assert(t, !DeepEqualStringArray(b, c), "Error DeepEqualArray")
	assert.Assert(t, !DeepEqualStringArray(b, d), "Error DeepEqualArray")
}
