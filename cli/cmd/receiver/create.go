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
	url := viper.GetString("url")
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

	name := viper.GetString("name")
	etype := viper.GetString("type")
	version := viper.GetString("version")
	desc := viper.GetString("description")
	schema := viper.GetString("schema")

	dryrun := viper.GetBool("dry-run")
	noindent := viper.GetBool("no-indent")

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
	createCmd.Flags().String("name", "", "Name of the Event Receiver")
	createCmd.Flags().String("type", "", "Type of the Event Receiver Group")
	createCmd.Flags().String("version", "", "Version of the Event Receiver Group")
	createCmd.Flags().String("description", "", "Description of the Event Receiver")
	createCmd.Flags().String("schema", "{}", "Schema of the Event Receiver")
	createCmd.Flags().String("url", "http://localhost:8042", "EPR base url")
	createCmd.Flags().Bool("dry-run", false, "do a dry run of the command")
	createCmd.Flags().Bool("no-indent", false, "do not indent the JSON output")

	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("type")
	_ = createCmd.MarkFlagRequired("version")
	_ = createCmd.MarkFlagRequired("description")
	_ = createCmd.MarkFlagRequired("schema")

	return createCmd
}
