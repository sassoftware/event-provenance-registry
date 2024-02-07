// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"reflect"
	"testing"
)

func TestToJSON(t *testing.T) {
	// Test case 1: data already has JSON prefix
	data := []byte(`{"name":"John"}`)
	expected := []byte(`{"name":"John"}`)
	result, err := ToJSON(data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Test case 2: data does not have JSON prefix
	data = []byte(`name: John`)
	expected = []byte(`{"name":"John"}`)
	result, err = ToJSON(data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
