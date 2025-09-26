// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package receiver

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
	searchIdFlag       = "id"
	searchNameFlag     = "name"
	searchVersionFlag  = "version"
	searchTypeFlag     = "type"
	searchFieldsFlag   = "fields"
	searchJsonpathFlag = "jsonpath"
	searchUrlFlag      = "url"
	searchDryRunFlag   = "dry-run"
	searchNoIndentFlag = "no-indent"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "Search for Event Receivers",
	Long:    `Search for Event Receivers`,
	PreRunE: common.BindFlagsE,
	RunE:    runEventReceiverSearch,
}

func runEventReceiverSearch(_ *cobra.Command, _ []string) error {
	dryrun := viper.GetBool(searchDryRunFlag)
	noindent := viper.GetBool(searchNoIndentFlag)

	url := viper.GetString(searchUrlFlag)
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

	params := make(map[string]interface{})

	id := viper.GetString(searchIdFlag)
	if id != "" {
		params["id"] = id
	}

	name := viper.GetString(searchNameFlag)
	if name != "" {
		params["name"] = name
	}

	version := viper.GetString(searchVersionFlag)
	if version != "" {
		params["version"] = version
	}

	typeStr := viper.GetString(searchTypeFlag)
	if typeStr != "" {
		params["type"] = typeStr
	}

	fields, err := common.ProcessSearchFields(
		viper.GetStringSlice(searchFieldsFlag),
		&storage.EventReceiver{},
	)
	if err != nil {
		return err
	}

	if dryrun {
		fmt.Printf("ID: %s\n", id)
		fmt.Printf("Name: %s\n", name)
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Type: %s\n", typeStr)
		fmt.Printf("Fields: %v\n", fields)
		curlcmd, err := c.GetCurlSearch("events", params, fields)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", curlcmd)
		return nil
	}

	eventReceivers, err := c.SearchEventReceivers(params, fields)
	if err != nil {
		return err
	}

	if noindent {
		content, err := json.Marshal(eventReceivers)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", content)
		return nil
	}
	content, err := json.MarshalIndent(eventReceivers, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", content)
	return nil
}

// NewSearchCmd returns a new search command
func NewSearchCmd() *cobra.Command {
	searchCmd.Flags().String(searchIdFlag, "", "ID for the event receiver")
	searchCmd.Flags().String(searchNameFlag, "", "Name of the event receiver")
	searchCmd.Flags().String(searchVersionFlag, "", "Version of the event receiver")
	searchCmd.Flags().String(searchTypeFlag, "", "Type of the event receiver")
	searchCmd.Flags().
		String(searchFieldsFlag, "id name version type", "Space delimited list of fields, or 'all' for all user fields")
	searchCmd.Flags().String(searchJsonpathFlag, "", "JSONPath expression to apply to output")
	searchCmd.Flags().String(searchUrlFlag, "http://localhost:8042", "EPR base url")
	searchCmd.Flags().Bool(searchDryRunFlag, false, "do a dry run of the command")
	searchCmd.Flags().Bool(searchNoIndentFlag, false, "do not indent the JSON output")
	return searchCmd
}
