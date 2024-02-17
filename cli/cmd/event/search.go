// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package event

import (
	"encoding/json"
	"fmt"

	"github.com/sassoftware/event-provenance-registry/cli/cmd/common"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "Searches for event objects",
	Long:    `Searches for event objects`,
	PreRunE: common.BindFlagsE,
	RunE:    runSearchEvent,
}

// runSearchEvent runs the search and return error
func runSearchEvent(_ *cobra.Command, _ []string) error {
	dryrun := viper.GetBool("dry-run")
	noindent := viper.GetBool("no-indent")

	url := viper.GetString("url")
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

	params := make(map[string]interface{})

	id := viper.GetString("id")
	if id != "" {
		params["id"] = id
	}

	name := viper.GetString("name")
	if name != "" {
		params["name"] = name
	}

	version := viper.GetString("version")
	if version != "" {
		params["version"] = version
	}

	release := viper.GetString("release")
	if release != "" {
		params["release"] = release
	}

	platformID := viper.GetString("platform_id")
	if platformID != "" {
		params["platform_id"] = platformID
	}

	pkg := viper.GetString("package")
	if pkg != "" {
		params["package"] = pkg
	}

	success := viper.GetString("success")
	if success != "" {
		params["success"] = success
	}

	eventReceiverID := viper.GetString("event-receiver-id")
	if eventReceiverID != "" {
		params["event_receiver_id"] = eventReceiverID
	}

	fields, err := common.ProcessSearchFields(viper.GetStringSlice("fields"), &storage.Event{})
	if err != nil {
		return err
	}

	if dryrun {
		fmt.Printf("ID: %s\n", id)
		fmt.Printf("Name: %s\n", name)
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Release: %s\n", release)
		fmt.Printf("PlatformID: %s\n", platformID)
		fmt.Printf("Package: %s\n", pkg)
		fmt.Printf("Success: %s\n", success)
		fmt.Printf("EventReceiverID: %s\n", eventReceiverID)
		fmt.Printf("Fields: %v\n", fields)
		curlcmd, err := c.GetCurlSearch("events", params, fields)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", curlcmd)
		return nil
	}

	events, err := c.SearchEvents(params, fields)
	if err != nil {
		return err
	}

	if noindent {
		content, err := json.Marshal(events)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", content)
		return nil
	}
	content, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", content)
	return nil
}

// NewSearchCmd returns a new search command
func NewSearchCmd() *cobra.Command {
	searchCmd.Flags().String("id", "", "Id for the event")
	searchCmd.Flags().String("name", "", "Name of the event")
	searchCmd.Flags().String("version", "", "Version of the event")
	searchCmd.Flags().String("release", "", "Release of the event")
	searchCmd.Flags().String("platform_id", "", "Platform id of the event")
	searchCmd.Flags().String("package", "", "Package of the event")
	searchCmd.Flags().String("success", "", "Success of the event")
	searchCmd.Flags().String("event-receiver-id", "", "Event receiver id of the event")
	searchCmd.Flags().String("fields", "id name version release platform_id package success", "Space delimited list of fields, or 'all' for all user fields")
	searchCmd.Flags().String("jsonpath", "", "JSONPath expression to apply to output")
	searchCmd.Flags().String("url", "http://localhost:8042", "EPR base url")
	searchCmd.Flags().Bool("dry-run", false, "do a dry run of the command")
	searchCmd.Flags().Bool("no-indent", false, "do not indent the JSON output")
	return searchCmd
}
