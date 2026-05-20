// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package event

import (
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/cli/cmd/common"
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flag constants
const (
	createNameFlag            = "name"
	createDescriptionFlag     = "description"
	createVersionFlag         = "version"
	createReleaseFlag         = "release"
	createPlatformIDFlag      = "platform-id"
	createPackageFlag         = "package"
	createSuccessFlag         = "success"
	createEventReceiverIDFlag = "event-receiver-id"
	createPayloadFlag         = "payload"
	createUrlFlag             = "url"
	createDryRunFlag          = "dry-run"
	createNoIndentFlag        = "no-indent"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "create an event",
	Long:    `create an event`,
	PreRunE: common.BindFlagsE,
	RunE:    runCreateEvent,
}

func runCreateEvent(_ *cobra.Command, _ []string) error {
	url := viper.GetString(createUrlFlag)
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

	name := viper.GetString(createNameFlag)
	description := viper.GetString(createDescriptionFlag)
	version := viper.GetString(createVersionFlag)
	release := viper.GetString(createReleaseFlag)
	platform := viper.GetString(createPlatformIDFlag)
	pkg := viper.GetString(createPackageFlag)
	success := viper.GetBool(createSuccessFlag)
	eventReceiverID := viper.GetString(createEventReceiverIDFlag)
	payload := viper.GetString(createPayloadFlag)
	dryrun := viper.GetBool(createDryRunFlag)
	noindent := viper.GetBool(createNoIndentFlag)

	e := &storage.Event{
		Name:            name,
		Description:     description,
		Version:         version,
		Release:         release,
		PlatformID:      platform,
		Package:         pkg,
		Success:         success,
		EventReceiverID: graphql.ID(eventReceiverID),
		Payload:         types.JSON{JSON: []byte(payload)},
	}

	if dryrun {
		content, err := e.ToJSON()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", content)
		return nil
	}

	content, err := c.CreateEvent(e)
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
	// createCmd.Flags().StringP("payload", "p", "", "Takes a JSON payload describing a event.")
	createCmd.Flags().String(createNameFlag, "", "name")
	createCmd.Flags().String(createDescriptionFlag, "", "description")
	createCmd.Flags().String(createVersionFlag, "", "version")
	createCmd.Flags().String(createReleaseFlag, "", "release string")
	createCmd.Flags().String(createPlatformIDFlag, "", "platform ID code")
	createCmd.Flags().String(createPackageFlag, "", "package type")
	createCmd.Flags().Bool(createSuccessFlag, false, "specify if the event succeeded")
	createCmd.Flags().String(createEventReceiverIDFlag, "", "ID of the event receiver")
	createCmd.Flags().String(createPayloadFlag, "", "JSON string of event payload")
	createCmd.Flags().String(createUrlFlag, "http://localhost:8042", "EPR base url")
	createCmd.Flags().Bool(createDryRunFlag, false, "do a dry run of the command")
	createCmd.Flags().Bool(createNoIndentFlag, false, "do not indent the JSON output")
	_ = createCmd.MarkFlagRequired(createNameFlag)
	_ = createCmd.MarkFlagRequired(createDescriptionFlag)
	_ = createCmd.MarkFlagRequired(createVersionFlag)
	_ = createCmd.MarkFlagRequired(createReleaseFlag)
	_ = createCmd.MarkFlagRequired(createPlatformIDFlag)
	_ = createCmd.MarkFlagRequired(createPackageFlag)
	_ = createCmd.MarkFlagRequired(createEventReceiverIDFlag)
	_ = createCmd.MarkFlagRequired(createPayloadFlag)

	return createCmd
}
