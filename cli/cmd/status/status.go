// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"fmt"
	"net/url"

	"github.com/sassoftware/event-provenance-registry/cli/cmd/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// gateCmd represents the gate command
var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Check status of EPR",
	Long:    `Check the status of EPR including liveness and readiness`,
	PreRunE: common.BindFlagsE,
	RunE:    run,
}

// run runs the call to create a status check, returns error
func run(_ *cobra.Command, _ []string) error {
	dryrun := viper.GetBool("dry-run")
	noindent := viper.GetBool("no-indent")

	u := viper.GetString("url")

	c, err := common.GetClient(u)
	if err != nil {
		return err
	}

	if dryrun {
		endpoint, err := url.JoinPath(u, "healthz", "status")
		if err != nil {
			return err
		}
		fmt.Printf("\ncurl -X GET %s\n", endpoint)
		return nil
	}
	status, err := c.CheckStatus()
	if err != nil {
		return err
	}

	if noindent {
		fmt.Printf("%s\n", status)
		return nil
	}

	content, err := common.IndentJSON(status)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", content)
	return nil
}

// NewStatusCmd returns the statusCmd
func NewStatusCmd() *cobra.Command {
	statusCmd.Flags().String("url", "http://localhost:8042", "EPR base url")
	statusCmd.Flags().Bool("dry-run", false, "do a dry run of the command")
	statusCmd.Flags().Bool("no-indent", false, "do not indent the JSON output")
	return statusCmd
}
