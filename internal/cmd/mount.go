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
	dbFactory "github.com/kos-v/dbunderfs/internal/factory/db"
	"github.com/kos-v/dbunderfs/internal/fs"
	"github.com/kos-v/dsnparser"
	"github.com/spf13/cobra"
)

type mountOpts struct {
	dsn   string
	point string
}

func mountCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "mount POINT DSN",
		Short: "Mounts database FS to a specified mount point",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMount(mountOpts{
				dsn:   args[0],
				point: args[1],
			})
		},
	}

	return command
}

func runMount(opts mountOpts) error {
	dsn := dsnparser.Parse(opts.dsn)

	dbInstance, err := dbFactory.CreateInstance(*dsn)
	if err != nil {
		return err
	}
	defer dbInstance.Close()

	repositoryRegistry, err := dbFactory.CreateRepositoryRegistry(dbInstance)
	if err != nil {
		return err
	}

	return fs.Mount(opts.point, repositoryRegistry)
}
