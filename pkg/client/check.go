// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"log/slog"
)

// CheckReadiness checks EPR readiness
func (c *Client) CheckReadiness() (bool, error) {
	endpoint, err := c.getHealthEndpoint("/readiness")
	if err != nil {
		return false, err
	}
	content, err := c.DoGet(endpoint)
	slog.Debug("Check Readiness", "content", content)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckLiveness checks the EPRs liveness
func (c *Client) CheckLiveness() (bool, error) {
	endpoint, err := c.getHealthEndpoint("/liveness")
	if err != nil {
		return false, err
	}
	content, err := c.DoGet(endpoint)
	slog.Debug("Check Liveness", "content", content)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckStatus checks the EPRs Status
func (c *Client) CheckStatus() (string, error) {
	endpoint, err := c.getHealthEndpoint("/status")
	if err != nil {
		return "", err
	}
	content, err := c.DoGet(endpoint)
	slog.Debug("Check Status", "content", content)
	if err != nil {
		return content, err
	}
	return content, nil
}
