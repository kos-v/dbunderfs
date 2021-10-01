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
	"github.com/kos-v/dbunderfs/src/db"
	"github.com/kos-v/dbunderfs/src/db/migration"
	dbFactory "github.com/kos-v/dbunderfs/src/factory/db"
	log "github.com/kos-v/dbunderfs/src/log"
	"github.com/kos-v/dsnparser"
	"github.com/spf13/cobra"
)

type migrateOpts struct {
	downDirectionNumber int
	direction           string
	dsn                 string
}

func migrateCommand() *cobra.Command {
	opts := migrateOpts{}
	command := &cobra.Command{
		Use:   "migrate DIRECTION DSN",
		Short: "Performs migration operations on a file system database",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] != migration.DirUp && args[0] != migration.DirDown {
				return fmt.Errorf("%q is unknown operation of migrate command. Should be %q or %q.", args[0], migration.DirUp, migration.DirDown)
			}

			opts.direction = args[0]
			opts.dsn = args[1]

			return runMigrate(opts)
		},
	}

	command.Flags().IntVar(&opts.downDirectionNumber, "downNumber", 1, "Number of migrations to rollback")

	return command
}

func runMigrate(opts migrateOpts) error {
	dsn := dsnparser.Parse(opts.dsn)
	dbInstance, err := dbFactory.CreateInstance(*dsn)
	if err != nil {
		return err
	}
	defer dbInstance.Close()

	migrator, err := createMigrator(dbInstance, opts)
	if err != nil {
		return err
	}

	err = migrator.Migrate()
	if err != nil {
		return err
	}

	return nil
}

func createMigrator(dbInstance db.DBInstance, opts migrateOpts) (*migration.Migrator, error) {
	commiter, err := dbFactory.CreateMigrationCommiter(dbInstance)
	if err != nil {
		return nil, err
	}

	migrations, err := dbFactory.CreateMigrations(dbInstance)
	if err != nil {
		return nil, err
	}

	return &migration.Migrator{
		Commiter:            commiter,
		DownDirectionNumber: opts.downDirectionNumber,
		Direction:           opts.direction,
		QueryExecutor:       dbInstance,
		Logger:              log.NewStdoutLogger(),
		Migrations:          migrations,
	}, nil
}
