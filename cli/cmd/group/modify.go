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

// Flag constants
const (
	modifyIdFlag        = "id"
	modifyDisableFlag   = "disable"
	modifyEnableFlag    = "enable"
	modifyUrlFlag       = "url"
	modifyDryRunFlag    = "dry-run"
	modifyNoIndentFlag  = "no-indent"
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
	id := viper.GetString(modifyIdFlag)
	disable := viper.GetBool(modifyDisableFlag)
	enable := viper.GetBool(modifyEnableFlag)
	dryrun := viper.GetBool(modifyDryRunFlag)
	noindent := viper.GetBool(modifyNoIndentFlag)

	url := viper.GetString(modifyUrlFlag)
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
	modifyCmd.Flags().String(modifyIdFlag, "", "ID of the Event Receiver Group")
	modifyCmd.Flags().Bool(modifyDisableFlag, false, "Disable the Event Receiver Group")
	modifyCmd.Flags().Bool(modifyEnableFlag, false, "Enable the Event Receiver Group")
	modifyCmd.Flags().String(modifyUrlFlag, "http://localhost:8042", "EPR base url")
	modifyCmd.Flags().Bool(modifyDryRunFlag, false, "do a dry run of the command")
	modifyCmd.Flags().Bool(modifyNoIndentFlag, false, "do not indent the JSON output")

	_ = modifyCmd.MarkFlagRequired(modifyIdFlag)

	return modifyCmd
}
