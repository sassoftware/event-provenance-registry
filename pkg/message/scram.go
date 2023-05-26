// Copyright © 2019, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package message

import (
	"crypto/sha512"

	"github.com/xdg/scram"
)

// SHA512 references the sha512 hash function.
var SHA512 scram.HashGeneratorFcn = sha512.New

// SCRAMClient client for doing SCRAM authentication through sarama. This implementation was taken from the sarama examples
// here: https://github.com/Shopify/sarama/blob/master/examples/sasl_scram_client/scram_client.go
type SCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

// Begin begins SCRAM authentication
func (s *SCRAMClient) Begin(user, pass, authzID string) error {
	c, err := s.HashGeneratorFcn.NewClient(user, pass, authzID)
	if err != nil {
		return err
	}
	s.Client = c
	s.ClientConversation = s.Client.NewConversation()
	return nil
}

// Step handles scram auth
func (s *SCRAMClient) Step(challenge string) (string, error) {
	return s.ClientConversation.Step(challenge)
}

// Done closes conversation
func (s *SCRAMClient) Done() bool {
	return s.ClientConversation.Done()
}
