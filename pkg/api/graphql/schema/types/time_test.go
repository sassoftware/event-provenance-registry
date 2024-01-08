// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"gorm.io/datatypes"
	"gotest.tools/v3/assert"
)

func TestUnmarshalGraphQL(t *testing.T) {
	// Positive test case
	input := datatypes.Date(time.Now())
	expectedOutput := &Time{Date: input}

	fmt.Printf("%v", input)
	time := &Time{}
	err := time.UnmarshalGraphQL(input)
	assert.NilError(t, err)

	if !reflect.DeepEqual(time, expectedOutput) {
		t.Errorf("expected %+v, but got %+v", expectedOutput, time)
	}

	// Negative test case
	invalidInput := "invalid"
	err = time.UnmarshalGraphQL(invalidInput)
	assert.Assert(t, nil != err)
	expectedError := fmt.Errorf("wrong type for Time: %T", invalidInput)
	assert.Equal(t, expectedError.Error(), err.Error(), "expected error message does not match results")
}

func TestMarshalJSON(t *testing.T) {
	// Positive test case
	expected := []byte(`"2022-01-01T00:00:00Z"`)
	date := datatypes.Date(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC))
	tm := &Time{Date: date}
	result, err := tm.MarshalJSON()
	assert.NilError(t, err)
	assert.Equal(t, string(result), string(expected), "expected %s, but got %s", expected, result)
}
