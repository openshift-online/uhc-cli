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

package describe

import (
	"fmt"
	"os"

	"github.com/ALimaRedHat/uhc-cli/pkg/dump"
	"github.com/openshift-online/uhc-cli/pkg/config"
	"github.com/openshift-online/uhc-cli/pkg/util"
	"github.com/openshift-online/uhc-sdk-go/pkg/client"

	"github.com/spf13/cobra"
)

var args struct {
	debug bool
	out   bool
	short bool
}

var Cmd = &cobra.Command{
	Use:   "describe CLUSTERID [--output] [--short]",
	Short: "Decribe a cluster",
	Long:  "Get info about a cluster identified by its cluster ID, optional output flag",
	Run:   run,
}

// Check error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	// Add flags to rootCmd
	flags := Cmd.Flags()
	flags.BoolVar(
		&args.debug,
		"debug",
		false,
		"Enable debug mode.",
	)
	flags.BoolVar(
		&args.out,
		"out",
		true,
		"Output result into json file.",
	)
	flags.BoolVar(
		&args.short,
		"short",
		false,
		"Output a simpler description of the defined cluster",
	)
}

func run(cmd *cobra.Command, argv []string) {

	if len(argv) != 1 {
		fmt.Fprintf(os.Stderr, "Expected exactly one cluster id\n")
		os.Exit(1)
	}

	// Load the configuration file:
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't load config file: %v\n", err)
		os.Exit(1)
	}
	if cfg == nil {
		fmt.Fprintf(os.Stderr, "Not logged in, run the 'login' command\n")
		os.Exit(1)
	}

	// Check that the configuration has credentials or tokens that haven't have expired:
	armed, err := config.Armed(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't check if tokens have expired: %v\n", err)
		os.Exit(1)
	}
	if !armed {
		fmt.Fprintf(os.Stderr, "Tokens have expired, run the 'login' command\n")
		os.Exit(1)
	}

	// Create the connection:
	logger, err := util.NewLogger(args.debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create logger: %v\n", err)
		os.Exit(1)
	}

	// Create the connection, and remember to close it:
	connection, err := client.NewConnectionBuilder().
		Logger(logger).
		TokenURL(cfg.TokenURL).
		Client(cfg.ClientID, cfg.ClientSecret).
		Scopes(cfg.Scopes...).
		URL(cfg.URL).
		User(cfg.User, cfg.Password).
		Tokens(cfg.AccessToken, cfg.RefreshToken).
		Insecure(cfg.Insecure).
		Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create connection: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	// Get full Api Response (JSON)
	if !args.short {
		// Get path
		request := connection.Get().Path("/api/clusters_mgmt/v1/clusters/" + argv[0])

		// Send the request:
		responseGet, err := request.Send()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't send request: %v\n", err)
			os.Exit(1)
		}

		body := responseGet.Bytes()
		check(err)
		// Terminal or Json output
		if args.out {
			err = dump.Pretty(os.Stdout, body)
		} else {
			// Output to json
			myFile, err := os.Create("output.json")
			check(err)
			e := dump.Pretty(myFile, body)
			check(e)
		}
		// Get Api response using cluster class for short output
	} else {
		// Get the client for the resource that manages the collection of clusters:
		resource := connection.ClustersMgmt().V1().Clusters()

		// Get the resource that manages the cluster that we want to display:
		clusterResource := resource.Cluster(argv[0])

		// Retrieve the collection of clusters:
		response, err := clusterResource.Get().
			Send()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't retrieve clusters: %s", err)
			os.Exit(1)
		}

		cluster := response.Body()
		fmt.Printf("ID:       %s\n"+
			"Name:     %s.%s\n"+
			"Masters:  %d\n"+
			"Computes: %d\n"+
			"Region:   %s\n"+
			"Multi-az: %t\n",
			cluster.ID(),
			cluster.Name(),
			cluster.DNS().BaseDomain(),
			cluster.Nodes().Master(),
			cluster.Nodes().Compute(),
			cluster.Region().ID(),
			cluster.MultiAZ(),
		)
	}
}
