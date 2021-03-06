/*
   Copyright 2021 The DbunderFS Contributors.

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

package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

type buildVersionInfo struct {
	build         string
	buildDatetime string
	release       string
}

type buildInfo struct {
	binary  string
	debug   bool
	version buildVersionInfo
}

func (i *buildInfo) GetSummary() string {
	summary := i.version.release + ", build " + i.version.build + " from " + i.version.buildDatetime
	if i.debug {
		summary += ", debug"
	}
	return summary
}

func RootCommand() *cobra.Command {
	bi := buildInfo{
		binary: fBinary,
		debug:  false,
		version: buildVersionInfo{
			build:         fBuild,
			buildDatetime: fBuildDatetime,
			release:       fRelease,
		},
	}

	if fDebug == "true" {
		bi.debug = true
	}

	command := &cobra.Command{
		Use:              bi.binary,
		Short:            "DbunderFS",
		Version:          bi.GetSummary(),
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
		migrateCommand(),
	)

	return command
}
