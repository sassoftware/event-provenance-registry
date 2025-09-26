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
	createNameFlag             = "name"
	createTypeFlag             = "type"
	createVersionFlag          = "version"
	createDescriptionFlag      = "description"
	createEventReceiverIDsFlag = "event-receiver-ids"
	createEnabledFlag          = "enabled"
	createUrlFlag              = "url"
	createDryRunFlag           = "dry-run"
	createNoIndentFlag         = "no-indent"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "creates a Event Receiver Group",
	Long:    `creates a Event Receiver Group`,
	PreRunE: common.BindFlagsE,
	RunE:    runCreateEventReceiverGroup,
}

// runCreateEventReceiverGroup creates the Event Receiver Group, returns error
func runCreateEventReceiverGroup(_ *cobra.Command, _ []string) error {
	url := viper.GetString(createUrlFlag)
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

	name := viper.GetString(createNameFlag)
	etype := viper.GetString(createTypeFlag)
	version := viper.GetString(createVersionFlag)
	desc := viper.GetString(createDescriptionFlag)
	evrIDs := viper.GetStringSlice(createEventReceiverIDsFlag)
	enabled := viper.GetBool(createEnabledFlag)
	dryrun := viper.GetBool(createDryRunFlag)
	noindent := viper.GetBool(createNoIndentFlag)

	eventReceiverIDs := []graphql.ID{}
	for _, id := range evrIDs {
		eventReceiverIDs = append(eventReceiverIDs, graphql.ID(id))
	}

	erg := &storage.EventReceiverGroup{
		Name:             name,
		Type:             etype,
		Version:          version,
		Description:      desc,
		EventReceiverIDs: eventReceiverIDs,
		Enabled:          enabled,
	}

	if dryrun {
		content, err := json.MarshalIndent(erg, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", content)
		return nil
	}

	content, err := c.CreateEventReceiverGroup(erg)
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

// NewCreateCmd creates a new command
func NewCreateCmd() *cobra.Command {
	createCmd.Flags().String(createNameFlag, "", "Name of the Event Receiver Group")
	createCmd.Flags().String(createTypeFlag, "", "Type of the Event Receiver Group")
	createCmd.Flags().String(createVersionFlag, "", "Version of the Event Receiver Group")
	createCmd.Flags().String(createDescriptionFlag, "", "Description of the Event Receiver Group")
	createCmd.Flags().String(createEventReceiverIDsFlag, "", "Space delimited set of receiver ids")
	createCmd.Flags().Bool(createEnabledFlag, true, "Enable the Event Receiver Group")
	createCmd.Flags().String(createUrlFlag, "http://localhost:8042", "EPR base url")
	createCmd.Flags().Bool(createDryRunFlag, false, "do a dry run of the command")
	createCmd.Flags().Bool(createNoIndentFlag, false, "do not indent the JSON output")
	_ = createCmd.MarkFlagRequired(createNameFlag)
	_ = createCmd.MarkFlagRequired(createTypeFlag)
	_ = createCmd.MarkFlagRequired(createDescriptionFlag)
	_ = createCmd.MarkFlagRequired(createEventReceiverIDsFlag)

	return createCmd
}
