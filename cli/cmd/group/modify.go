// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package group

import (
	"encoding/json"
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/cli/cmd/common"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var modifyCmd = &cobra.Command{
	Use:     "modify",
	Short:   "Modifies a stage",
	Long:    `modifies a stage`,
	PreRunE: common.BindFlagsE,
	RunE:    runModifyEventReceiverGroup,
}

// runModifyEventReceiverGroup modifys stage and returns error
func runModifyEventReceiverGroup(_ *cobra.Command, _ []string) error {
	id := viper.GetString("id")
	disable := viper.GetBool("disable")
	enable := viper.GetBool("enable")
	dryrun := viper.GetBool("dry-run")
	noindent := viper.GetBool("no-indent")

	url := viper.GetString("url")
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}
	enabled := true
	if disable {
		enabled = false
	}
	if enable {
		enabled = true
	}
	erg := &storage.EventReceiverGroup{
		ID:      graphql.ID(id),
		Enabled: enabled,
	}

	if dryrun {
		content, err := json.Marshal(erg)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", content)
		return nil
	}

	content, err := c.ModifyEventReceiverGroup(erg)
	if err != nil {
		return err
	}

	if noindent {
		fmt.Printf("%s\n", content)
		return nil
	}

	content, err = common.IndentJSON(content)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", content)
	return nil
}

// NewModifyCmd creates a new command
func NewModifyCmd() *cobra.Command {
	modifyCmd.Flags().String("id", "", "ID of the Event Receiver Group")
	modifyCmd.Flags().Bool("disable", false, "Disable the Event Receiver Group")
	modifyCmd.Flags().Bool("enable", false, "Enable the Event Receiver Group")

	_ = modifyCmd.MarkFlagRequired("id")

	return modifyCmd
}
