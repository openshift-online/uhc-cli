package command

import "github.com/spf13/cobra"

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

// This file contains the types and functions used to manage the configuration of the command line
// client.

// StringInArray returns true if it finds the given string in the given array, else false
func StringInArray(str string, arr []string) bool {

	for _, element := range arr {
		if str == element {
			return true
		}
	}

	return false
}

// GetCommandNames takes a list of commands and returns a list of their names
func GetCommandNames(commandList []*cobra.Command) []string {
	var commandNames []string

	for _, command := range commandList {
		commandNames = append(commandNames, command.Name())
	}

	return commandNames
}
