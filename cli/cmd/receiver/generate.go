// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package receiver

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
	Short:   "Generates a Event Receiver",
	Long:    `Generates a Event Receiver from a file`,
	Args:    cobra.ExactArgs(1),
	PreRunE: common.BindFlagsE,
	RunE:    runGenerateEventReceiver,
}

// runGenerateEventReceiver will create a Event Receiver with the content based in a file, returns error
func runGenerateEventReceiver(_ *cobra.Command, args []string) error {
	dryrun := viper.GetBool("dry-run")
	url := viper.GetString("url")
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

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

	jsonStr, err := common.ToJSON(byteValue)
	if err != nil {
		return err
	}

	ers := []*storage.EventReceiver{}

	if common.IsJSONArray(jsonStr) {
		err = json.Unmarshal(jsonStr, &ers)
		if err != nil {
			return err
		}
	} else {
		_, err := jsonFile.Seek(0, 0)
		if err != nil {
			return err
		}
		er, err := storage.EventReceiverFromJSON(jsonFile)
		if err != nil {
			return err
		}
		ers = append(ers, er)
	}

	if dryrun {
		content, err := json.MarshalIndent(ers, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", content)
		return nil
	}

	for idx := range ers {
		content, err := c.CreateEventReceiver(ers[idx])
		if err != nil {
			return err
		}
		fmt.Println(content)
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
