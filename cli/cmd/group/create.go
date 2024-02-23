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
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "creates a Event Receiver Group",
	Long:    `creates a Event Receiver Group`,
	PreRunE: common.BindFlagsE,
	RunE:    runCreateEventReceiverGroup,
}

// runCreateEventReceiverGroup creates the Event Receiver Group, returns error
func runCreateEventReceiverGroup(_ *cobra.Command, _ []string) error {
	url := viper.GetString("url")
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

	name := viper.GetString("name")
	etype := viper.GetString("type")
	version := viper.GetString("version")
	desc := viper.GetString("description")
	evrIDs := viper.GetStringSlice("event-receiver-ids")
	enabled := viper.GetBool("enabled")
	dryrun := viper.GetBool("dry-run")
	noindent := viper.GetBool("no-indent")

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
	createCmd.Flags().String("name", "", "Name of the Event Receiver Group")
	createCmd.Flags().String("type", "", "Type of the Event Receiver Group")
	createCmd.Flags().String("version", "", "Version of the Event Receiver Group")
	createCmd.Flags().String("description", "", "Description of the Event Receiver Group")
	createCmd.Flags().String("event-receiver-ids", "", "Space delimited set of receiver ids")
	createCmd.Flags().Bool("enabled", true, "Enable the Event Receiver Group")
	createCmd.Flags().String("url", "http://localhost:8042", "EPR base url")
	createCmd.Flags().Bool("dry-run", false, "do a dry run of the command")
	createCmd.Flags().Bool("no-indent", false, "do not indent the JSON output")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("type")
	_ = createCmd.MarkFlagRequired("description")
	_ = createCmd.MarkFlagRequired("event-receiver-ids")

	return createCmd
}
