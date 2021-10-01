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

package db

import (
	"fmt"
	"github.com/kos-v/dbunderfs/src/db"
	"github.com/kos-v/dbunderfs/src/db/migration"
	mysqlMigrations "github.com/kos-v/dbunderfs/src/migrations/mysql"
	"github.com/kos-v/dsnparser"
)

type DriverNotFoundError struct{ driver string }

func (err *DriverNotFoundError) Error() string {
	return fmt.Sprintf("driver for database %q not found.", err.driver)
}

func CreateInstance(dsn dsnparser.DSN) (db.DBInstance, error) {
	switch dsn.GetScheme() {
	case "mysql":
		inst := db.MySQLInstance{DSN: db.MySQLDSN{ParsedDSN: dsn}}
		if _, err := inst.Connect(); err != nil {
			return nil, err
		}
		return &inst, nil
	}

	return nil, &DriverNotFoundError{driver: dsn.GetScheme()}
}

func CreateMigrationCommiter(instance db.DBInstance) (migration.Commiter, error) {
	switch instance.GetDriverName() {
	case "mysql":
		return &migration.MySQLCommiter{Instance: instance}, nil
	}

	return nil, &DriverNotFoundError{driver: instance.GetDriverName()}
}

func CreateMigrations(instance db.DBInstance) ([]*migration.Migration, error) {
	switch instance.GetDriverName() {
	case "mysql":
		return mysqlMigrations.Migrations(), nil
	}

	return nil, &DriverNotFoundError{driver: instance.GetDriverName()}
}