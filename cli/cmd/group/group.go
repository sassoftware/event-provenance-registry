// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package group

import (
	"github.com/spf13/cobra"
)

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Create, Search, or Generate Event Receiver Groups",
	Long:  `Create, Search, or Generate Event Receiver Groups for the Event Provenance Registry Service`,
}

// NewGroupCmd create a new command
func NewGroupCmd() *cobra.Command {
	searchCmd := NewSearchCmd()
	groupCmd.AddCommand(searchCmd)
	createCmd := NewCreateCmd()
	groupCmd.AddCommand(createCmd)
	generateCmd := NewGenerateCmd()
	groupCmd.AddCommand(generateCmd)

	return groupCmd
}
