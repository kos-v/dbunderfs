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
	"bazil.org/fuse"
	"github.com/spf13/cobra"
)

func unmountCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "unmount POINT",
		Short: "Unmounts mount point from the filesystem",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUnmount(args[0])
		},
	}

	return command
}

func runUnmount(point string) error {
	return fuse.Unmount(point)
}