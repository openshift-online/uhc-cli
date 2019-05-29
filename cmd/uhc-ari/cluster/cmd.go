/*
Copyright (c) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"fmt"
	"os"

	"github.com/ALimaRedHat/uhc-cli/cmd/uhc/cluster/count"
	"github.com/ALimaRedHat/uhc-cli/cmd/uhc/cluster/describe"
	"github.com/ALimaRedHat/uhc-cli/cmd/uhc/cluster/list"
	"github.com/ALimaRedHat/uhc-cli/cmd/uhc/cluster/status"
	"github.com/spf13/cobra"
)

var args struct {
	debug bool
}
var Cmd = &cobra.Command{
	Use:   "cluster COMMAND",
	Short: "Get information about clusters",
	Long:  "Get status or information about a single cluster, or a list of clusters",
	Run:   run,
}

func init() {
	flags := Cmd.Flags()
	flags.BoolVar(
		&args.debug,
		"debug",
		false,
		"Enable debug mode.",
	)
	Cmd.AddCommand(count.Cmd)
	Cmd.AddCommand(list.Cmd)
	Cmd.AddCommand(status.Cmd)
	Cmd.AddCommand(describe.Cmd)
}

func run(cmd *cobra.Command, argv []string) {
	// Check there is at least one argument
	if len(argv) < 1 {
		fmt.Fprintf(os.Stderr, "Expected at least one argument\n")
		os.Exit(1)
	}
}
