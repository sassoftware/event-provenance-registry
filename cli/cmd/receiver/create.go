// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package receiver

import (
	"fmt"

	"github.com/sassoftware/event-provenance-registry/cli/cmd/common"
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flag constants
const (
	createNameFlag        = "name"
	createTypeFlag        = "type"
	createVersionFlag     = "version"
	createDescriptionFlag = "description"
	createSchemaFlag      = "schema"
	createUrlFlag         = "url"
	createDryRunFlag      = "dry-run"
	createNoIndentFlag    = "no-indent"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Creates a Event Receiver",
	Long:    `Creates a Event Receiver`,
	PreRunE: common.BindFlagsE,
	RunE:    runCreateEventReceiver,
}

// runCreateEventReceiver runs the call to create a EventReceiver, returns error
func runCreateEventReceiver(_ *cobra.Command, _ []string) error {
	url := viper.GetString(createUrlFlag)
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

	name := viper.GetString(createNameFlag)
	etype := viper.GetString(createTypeFlag)
	version := viper.GetString(createVersionFlag)
	desc := viper.GetString(createDescriptionFlag)
	schema := viper.GetString(createSchemaFlag)

	dryrun := viper.GetBool(createDryRunFlag)
	noindent := viper.GetBool(createNoIndentFlag)

	er := &storage.EventReceiver{
		Name:        name,
		Type:        etype,
		Version:     version,
		Description: desc,
		Schema:      types.JSON{JSON: []byte(schema)},
	}

	if dryrun {
		content, err := er.ToJSON()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", content)
		return nil
	}

	content, err := c.CreateEventReceiver(er)
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

// NewCreateCmd creates a new cmdline
func NewCreateCmd() *cobra.Command {
	createCmd.Flags().String(createNameFlag, "", "Name of the Event Receiver")
	createCmd.Flags().String(createTypeFlag, "", "Type of the Event Receiver Group")
	createCmd.Flags().String(createVersionFlag, "", "Version of the Event Receiver Group")
	createCmd.Flags().String(createDescriptionFlag, "", "Description of the Event Receiver")
	createCmd.Flags().String(createSchemaFlag, "{}", "Schema of the Event Receiver")
	createCmd.Flags().String(createUrlFlag, "http://localhost:8042", "EPR base url")
	createCmd.Flags().Bool(createDryRunFlag, false, "do a dry run of the command")
	createCmd.Flags().Bool(createNoIndentFlag, false, "do not indent the JSON output")

	_ = createCmd.MarkFlagRequired(createNameFlag)
	_ = createCmd.MarkFlagRequired(createTypeFlag)
	_ = createCmd.MarkFlagRequired(createVersionFlag)
	_ = createCmd.MarkFlagRequired(createDescriptionFlag)
	_ = createCmd.MarkFlagRequired(createSchemaFlag)

	return createCmd
}
