// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package response

// RestResponse generic rest response for all object types.
type RestResponse struct {
	Data   any     `json:"data"`
	Errors []error `json:"errors"`
}
