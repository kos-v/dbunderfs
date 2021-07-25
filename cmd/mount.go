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
	"github.com/kos-v/dbunderfs/db"
	"github.com/kos-v/dbunderfs/fs"
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

	dbInstance, err := createDBInstance(*dsn)
	if err != nil {
		return err
	}
	defer dbInstance.Close()

	dbFactory, err := createDBFactory(dbInstance)
	if err != nil {
		return err
	}

	if err = fs.Mount(opts.point, dbFactory); err != nil {
		return err
	}

	return nil
}

func createDBInstance(dsn dsnparser.DSN) (db.DBInstance, error) {
	if dsn.GetScheme() != "mysql" {
		return nil, fmt.Errorf("Driver for database \"%s\" not found.\n", dsn.GetScheme())
	}

	protocol := "tcp"
	if dsn.HasParam("protocol") {
		protocol = dsn.GetParam("protocol")
	}
	driverDsn := dsn.GetUser() + ":" + dsn.GetPassword() + "@" + protocol + "(" + dsn.GetHost() + ")/" + dsn.GetPath()

	inst := db.MySQLInstance{DSN: driverDsn}
	if _, err := inst.Connect(); err != nil {
		return nil, err
	}

	return &inst, nil
}

func createDBFactory(inst db.DBInstance) (db.DBFactory, error) {
	if inst.GetDriverName() != "mysql" {
		return nil, fmt.Errorf("Driver for database \"%s\" not found.\n", inst.GetDriverName())
	}

	return &db.MySQLFactory{Instance: inst}, nil
}
