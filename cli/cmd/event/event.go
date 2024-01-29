// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package event

import (
	"github.com/spf13/cobra"
)

// eventCmd represents the event command
var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Create, Search, and Generate Events",
	Long:  `Create, Search, and Generate Events for the Event Provenance Registry Service`,
}

// NewEventCmd command for new events
func NewEventCmd() *cobra.Command {
	searchCmd := NewSearchCmd()
	eventCmd.AddCommand(searchCmd)
	createCmd := NewCreateCmd()
	eventCmd.AddCommand(createCmd)
	generateCmd := NewGenerateCmd()
	eventCmd.AddCommand(generateCmd)
	return eventCmd
}
