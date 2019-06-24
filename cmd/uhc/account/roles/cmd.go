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

package roles

import (
	"fmt"
	"os"

	"github.com/openshift-online/uhc-cli/pkg/config"
	"github.com/openshift-online/uhc-cli/pkg/dump"
	"github.com/openshift-online/uhc-cli/pkg/util"
	"github.com/openshift-online/uhc-sdk-go/pkg/client"
	amv1 "github.com/openshift-online/uhc-sdk-go/pkg/client/accountsmgmt/v1"

	"github.com/spf13/cobra"
)

var args struct {
	debug bool
}

var Cmd = &cobra.Command{
	Use:   "roles [role-name]",
	Short: "Retrieve information of the different roles",
	Long:  "Get description of a role or list of all roles ",
	Run:   run,
}

func init() {
	// Add flags to rootCmd:
	flags := Cmd.Flags()
	flags.BoolVar(
		&args.debug,
		"debug",
		false,
		"Enable debug mode.",
	)
}

func run(cmd *cobra.Command, argv []string) {

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

	// No role name was provided,
	// Print all roles.
	if len(argv) < 1 {
		rolesListRequest := connection.AccountsMgmt().V1().Roles().List()
		response, err := rolesListRequest.Send()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't send request: %v\n", err)
			os.Exit(1)
		}
		response.Items().Each(func(item *amv1.Role) bool {
			fmt.Println(item.ID())
			return true
		})
		// Role name was provided.
	} else {

		// Get role with provided id response:
		roleResponse, err := connection.AccountsMgmt().V1().Roles().Role(argv[0]).Get().
			Send()
		role := roleResponse.Body()

		// Use role in new get request since original provides incomplete data.
		byteRole, err := connection.Get().Path(role.HREF()).
			Send()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't send request: %v\n", err)
			os.Exit(1)
		}

		// Dump pretty:
		dump.Pretty(os.Stdout, byteRole.Bytes())

	}

}
