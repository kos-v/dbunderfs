/*
   Copyright 2021 The DbunderFS Authors.

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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var binary string
var build string
var buildDatetime string
var version string

func RootCommand() *cobra.Command {
	version = version + ", build " + build + " from " + buildDatetime
	command := &cobra.Command{
		Use:              binary,
		Short:            "DbunderFS",
		Version:          version,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return fmt.Errorf("Unknown command: %q", "dbfs "+args[0])
		},
	}

	command.AddCommand(
		mountCommand(),
		unmountCommand(),
	)

	return command
}
