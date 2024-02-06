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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "create an event",
	Long:    `create an event`,
	PreRunE: common.BindFlagsE,
	RunE:    runCreateEvent,
}

func runCreateEvent(_ *cobra.Command, _ []string) error {
	url := viper.GetString("url")
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

	name := viper.GetString("name")
	description := viper.GetString("description")
	version := viper.GetString("version")
	release := viper.GetString("release")
	platform := viper.GetString("platform-id")
	pkg := viper.GetString("package")
	success := viper.GetBool("success")
	eventReceiverID := viper.GetString("receiver-id")
	payload := viper.GetString("payload")
	dryrun := viper.GetBool("dry-run")
	noindent := viper.GetBool("no-indent")

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
	createCmd.Flags().String("name", "", "name")
	createCmd.Flags().String("description", "", "description")
	createCmd.Flags().String("version", "", "version")
	createCmd.Flags().String("release", "", "release string")
	createCmd.Flags().String("platform-id", "", "platform ID code")
	createCmd.Flags().String("package", "", "package type")
	createCmd.Flags().Bool("success", false, "specify if the event succeeded")
	createCmd.Flags().String("receiver-id", "", "ID of the event receiver")
	createCmd.Flags().String("payload", "", "JSON string of event payload")
	createCmd.Flags().String("url", "http://localhost:8042", "EPR base url")
	createCmd.Flags().Bool("dry-run", false, "do a dry run of the command")
	createCmd.Flags().Bool("no-indent", false, "do not indent the JSON output")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("description")
	_ = createCmd.MarkFlagRequired("version")
	_ = createCmd.MarkFlagRequired("release")
	_ = createCmd.MarkFlagRequired("platform-id")
	_ = createCmd.MarkFlagRequired("package")
	_ = createCmd.MarkFlagRequired("receiver-id")
	_ = createCmd.MarkFlagRequired("payload")

	return createCmd
}
