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

// Flag constants
const (
	idFlag              = "id"
	nameFlag            = "name"
	versionFlag         = "version"
	releaseFlag         = "release"
	platformIDFlag      = "platform-id"
	packageFlag         = "package"
	successFlag         = "success"
	eventReceiverIDFlag = "event-receiver-id"
	fieldsFlag          = "fields"
	jsonpathFlag        = "jsonpath"
	urlFlag             = "url"
	dryRunFlag          = "dry-run"
	noIndentFlag        = "no-indent"
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
	dryrun := viper.GetBool(dryRunFlag)
	noindent := viper.GetBool(noIndentFlag)

	url := viper.GetString(urlFlag)
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

	params := make(map[string]interface{})

	id := viper.GetString(idFlag)
	if id != "" {
		params["id"] = id
	}

	name := viper.GetString(nameFlag)
	if name != "" {
		params["name"] = name
	}

	version := viper.GetString(versionFlag)
	if version != "" {
		params["version"] = version
	}

	release := viper.GetString(releaseFlag)
	if release != "" {
		params["release"] = release
	}

	platformID := viper.GetString(platformIDFlag)
	if platformID != "" {
		params["platform_id"] = platformID
	}

	pkg := viper.GetString(packageFlag)
	if pkg != "" {
		params["package"] = pkg
	}

	success := viper.GetString(successFlag)
	if success != "" {
		params["success"] = success
	}

	eventReceiverID := viper.GetString(eventReceiverIDFlag)
	if eventReceiverID != "" {
		params["event_receiver_id"] = eventReceiverID
	}

	fields, err := common.ProcessSearchFields(viper.GetStringSlice(fieldsFlag), &storage.Event{})
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
	searchCmd.Flags().String(idFlag, "", "Id for the event")
	searchCmd.Flags().String(nameFlag, "", "Name of the event")
	searchCmd.Flags().String(versionFlag, "", "Version of the event")
	searchCmd.Flags().String(releaseFlag, "", "Release of the event")
	searchCmd.Flags().String(platformIDFlag, "", "Platform id of the event")
	searchCmd.Flags().String(packageFlag, "", "Package of the event")
	searchCmd.Flags().String(successFlag, "", "Success of the event")
	searchCmd.Flags().String(eventReceiverIDFlag, "", "Event receiver id of the event")
	searchCmd.Flags().
		String(fieldsFlag, "id name version release platform_id package success", "Space delimited list of fields, or 'all' for all user fields")
	searchCmd.Flags().String(jsonpathFlag, "", "JSONPath expression to apply to output")
	searchCmd.Flags().String(urlFlag, "http://localhost:8042", "EPR base url")
	searchCmd.Flags().Bool(dryRunFlag, false, "do a dry run of the command")
	searchCmd.Flags().Bool(noIndentFlag, false, "do not indent the JSON output")
	return searchCmd
}
