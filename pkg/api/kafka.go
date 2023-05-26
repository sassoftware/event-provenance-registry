// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"sync"

	"gitlab.sas.com/async-event-infrastructure/server/pkg/config"
)

// SetupKafkaV1 sets up the kafka producer for v1
func SetupKafkaV1(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-cfg.Kafka.MsgChannel:
				cfg.Kafka.Producer.Async(cfg.Kafka.Topic, msg)
			}
		}
	}()
}
