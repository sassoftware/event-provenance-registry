// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package receiver

import (
	"github.com/sassoftware/event-provenance-registry/pkg/utils"
	"github.com/spf13/cobra"
)

var logger = utils.MustGetLogger("client", "cli.cmd.receiver")

// gateCmd represents the gate command
var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "Create, Search, or Generate Event Receivers",
	Long: `Create, Search, or Generate Event Receivers
	for the Event Provenance Registry Service`,
}

// NewReceiverCmd returns the receiverCmd
func NewReceiverCmd() *cobra.Command {
	searchCmd := NewSearchCmd()
	receiverCmd.AddCommand(searchCmd)
	createCmd := NewCreateCmd()
	receiverCmd.AddCommand(createCmd)
	generateCmd := NewGenerateCmd()
	receiverCmd.AddCommand(generateCmd)
	return receiverCmd
}
