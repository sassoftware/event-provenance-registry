// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

// CheckReadiness checks EPR readiness
func (c *Client) CheckReadiness() (bool, error) {
	endpoint := c.getHealthEndpoint("/readiness")
	content, err := c.DoGet(endpoint)
	logger.V(1).Info("Check Readiness : %s\n", content)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckLiveness checks the EPRs liveness
func (c *Client) CheckLiveness() (bool, error) {
	endpoint := c.getHealthEndpoint("/liveness")
	content, err := c.DoGet(endpoint)
	logger.V(1).Info("Check Liveness : %s\n", content)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckStatus checks the EPRs Status
func (c *Client) CheckStatus() (string, error) {
	endpoint := c.getHealthEndpoint("/status")
	content, err := c.DoGet(endpoint)
	logger.V(1).Info("Check Status : %s\n", content)
	if err != nil {
		return content, err
	}
	return content, nil
}
