// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package event

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sassoftware/event-provenance-registry/cli/cmd/common"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// generateCmd represents the create command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generates a event",
	Long:    `Generates a event from a file`,
	Args:    cobra.ExactArgs(1),
	PreRunE: common.BindFlagsE,
	RunE:    runGenerateEvent,
}

// runGenerateEvent runs the main command function
func runGenerateEvent(_ *cobra.Command, args []string) error {
	url := viper.GetString("url")
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

	dryrun := viper.GetBool("dry-run")
	noindent := viper.GetBool("no-indent")

	file := args[0]
	jsonFile, err := os.Open(file)
	if err != nil {
		jsonFile.Close()
		return fmt.Errorf("error opening file")
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("error reading file")
	}

	events := []*storage.Event{}
	jsonStr, err := common.ToJSON(byteValue)
	if err != nil {
		return err
	}
	if common.IsJSONArray(jsonStr) {
		err = json.Unmarshal(jsonStr, &events)
		if err != nil {
			return err
		}
	} else {
		_, err := jsonFile.Seek(0, 0)
		if err != nil {
			return err
		}
		ev, err := storage.EventFromJSON(jsonFile)
		if err != nil {
			return err
		}
		events = append(events, ev)
	}

	if dryrun {
		content, err := json.Marshal(events)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", content)
		return nil
	}

	for idx := range events {
		content, err := c.CreateEvent(events[idx])
		if err != nil {
			return err
		}
		if noindent {
			fmt.Printf("%s\n", content)
			continue
		}

		content, err = common.IndentJSON(content)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", content)
	}

	return nil
}

// NewGenerateCmd creates a new cmdline
func NewGenerateCmd() *cobra.Command {
	generateCmd.Flags().String("url", "http://localhost:8042", "EPR base url")
	generateCmd.Flags().Bool("dry-run", false, "do a dry run of the command")
	generateCmd.Flags().Bool("no-indent", false, "do not indent the JSON output")
	return generateCmd
}
