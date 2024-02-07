// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package group

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
	Short:   "Generates a Event Receiver Group",
	Long:    `Generates a Event Receiver Group from a file`,
	Args:    cobra.ExactArgs(1),
	PreRunE: common.BindFlagsE,
	RunE:    runGenerateEventReceiverGroup,
}

// runGenerateEventReceiverGroup creates Event Receiver Group from file contents, returns error
func runGenerateEventReceiverGroup(_ *cobra.Command, args []string) error {
	dryrun := viper.GetBool("dry-run")
	noindent := viper.GetBool("no-indent")

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

	ergs := []*storage.EventReceiverGroup{}
	jsonStr, err := common.ToJSON(byteValue)
	if err != nil {
		return err
	}
	if common.IsJSONArray(jsonStr) {
		err = json.Unmarshal(jsonStr, &ergs)
		if err != nil {
			return err
		}
	} else {
		_, err := jsonFile.Seek(0, 0)
		if err != nil {
			return err
		}
		erg, err := storage.EventReceiverGroupFromJSON(jsonFile)
		if err != nil {
			return err
		}
		ergs = append(ergs, erg)
	}

	if dryrun {
		content, err := json.MarshalIndent(ergs, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", content)
		return nil
	}

	for idx := range ergs {
		content, err := c.CreateEventReceiverGroup(ergs[idx])
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
