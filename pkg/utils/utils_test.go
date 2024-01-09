// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetTimestampRFC3339(t *testing.T) {
	input := "01F8TNE3GVBKDMZW9YWKYXX2N1"
	expected := "2021-06-22T16:08:55-04:00"
	res, err := GetTimestampRFC3339(input)
	assert.NilError(t, err)
	assert.Equal(t, expected, res)
	invalid := "ASKJDLKASD"
	_, errr := GetTimestampRFC3339(invalid)
	assert.Error(t, errr, "ulid: bad data size when unmarshaling")
}

func TestGetEnv(t *testing.T) {
	os.Setenv("FOO", "1")
	foo := GetEnv("FOO", "2")
	assert.Assert(t, foo == "1")
	bar := GetEnv("BAR", "42")
	assert.Assert(t, bar == "42")
}

func TestNewULID(t *testing.T) {
	u, err := NewULID()
	assert.NilError(t, err, "error is not nil")
	assert.Assert(t, len(u.String()) == 26)
}

func TestNewULIDAsString(t *testing.T) {
	u := NewULIDAsString()
	assert.Assert(t, len(u) == 26)
}

func TestNewUUID(t *testing.T) {
	u, err := NewUUID()
	assert.NilError(t, err, "error is not nil")
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
	assert.NilError(t, err, "error is not nil")
	equal, err := bc.Equal(hashedAPIKey)
	assert.NilError(t, err, "error is not nil")
	assert.Assert(t, equal)
}

// TestNewBCryptWithDifferentApiKeys validates two different api key hashs aren't equal
func TestNewBCryptWithDifferentApiKeys(t *testing.T) {
	bc := NewBCrypt("def456")
	hashedAPIKey, err := bc.Hash()
	assert.NilError(t, err, "error is not nil")
	_ = hashedAPIKey
	bc2 := NewBCrypt("jlp456")
	hashedAPIKey2, err := bc2.Hash()
	assert.NilError(t, err, "error is not nil")
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

func TestMustGetLogger(t *testing.T) {
	var logger = MustGetLogger("utils", "test.test")

	logger.Info("testing the info logger")
	logger.V(1).Info("testing the debug logger")
	logger.Error(fmt.Errorf("this is an error %s", "foobar"), "testing the error logger", "user", "foo")
	assert.Assert(t, logger != nil, "Error Logger is nil")
}
