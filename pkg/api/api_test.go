// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"testing"

	openapi "github.com/nasa9084/go-openapi"
)

func TestOpenAPIValidate(t *testing.T) {
	doc, err := openapi.LoadFile("../../resources/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if doc.Version != "3.1.0" {
		t.Error("OpenAPI version incorrect")
	}
}
