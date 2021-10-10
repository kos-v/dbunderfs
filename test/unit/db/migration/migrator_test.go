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

package migration

import (
	"github.com/kos-v/dbunderfs/internal/container"
	"github.com/kos-v/dbunderfs/internal/db/migration"
	dbTestHelpers "github.com/kos-v/dbunderfs/test/helpers/db"
	"github.com/kos-v/dbunderfs/test/helpers/factory/db"
	"testing"
)

func TestMigrator_Migrate_QueryResult(t *testing.T) {
	tests := []struct {
		direction       string
		downNumber      int
		migrations      []*migration.Migration
		upMigrations    []string
		expectedQueries []string
	}{
		{
			direction:    migration.DirUp,
			downNumber:   1,
			migrations:   testMigrations(),
			upMigrations: []string{},
			expectedQueries: []string{
				"UPDATE table1 SET col1 = \"val\"",
				"UPDATE table1 SET col2 = \"val\"",
				"UPDATE table2 SET col1 = \"val\"",
				"UPDATE table3 SET col1 = \"val\"",
			},
		},
		{
			direction:       migration.DirDown,
			downNumber:      -1,
			migrations:      testMigrations(),
			upMigrations:    []string{},
			expectedQueries: []string{},
		},
		{
			direction:    migration.DirUp,
			downNumber:   1,
			migrations:   testMigrations(),
			upMigrations: []string{"000000000000", "202102101045", "202103172206"},
			expectedQueries: []string{
				"UPDATE table3 SET col1 = \"val\"",
			},
		},
		{
			direction:    migration.DirDown,
			downNumber:   -1,
			migrations:   testMigrations(),
			upMigrations: []string{"000000000000", "202102101045", "202103172206"},
			expectedQueries: []string{
				"UPDATE table2 SET col1 = \"old val\"",
				"UPDATE table1 SET col2 = \"old val\"",
				"UPDATE table1 SET col1 = \"old val\"",
			},
		},
		{
			direction:       migration.DirUp,
			downNumber:      1,
			migrations:      testMigrations(),
			upMigrations:    []string{"000000000000", "202102101045", "202103172206", "202102150011"},
			expectedQueries: []string{},
		},
		{
			direction:    migration.DirDown,
			downNumber:   -1,
			migrations:   testMigrations(),
			upMigrations: []string{"000000000000", "202102101045", "202103172206", "202102150011"},
			expectedQueries: []string{
				"UPDATE table3 SET col1 = \"old val\"",
				"UPDATE table2 SET col1 = \"old val\"",
				"UPDATE table1 SET col2 = \"old val\"",
				"UPDATE table1 SET col1 = \"old val\"",
			},
		},
	}

	for testId, test := range tests {
		queryCollector := &dbTestHelpers.QueryExecuteCollector{Queries: &container.Collection{}}
		migrator := db.CreateMigrator(test.direction, test.downNumber, test.migrations, test.upMigrations)
		migrator.QueryExecutor = queryCollector

		if err := migrator.Migrate(); err != nil {
			t.Fatalf("Test %v fail: method Migrate returned an unexpected error. Error: %s", testId, err.Error())
		}

		if queryCollector.Queries.Len() != len(test.expectedQueries) {
			t.Fatalf("Test %v fail: the number of items is not as expected.\nExpected: %v. Result: %v.\n", testId, len(test.expectedQueries), queryCollector.Queries.Len())
		}

		for i, resultQuery := range queryCollector.Queries.ToList() {
			if resultQuery != test.expectedQueries[i] {
				t.Fatalf("Test %v fail: result data is not as expected.\nExpected: %v. Result: %v.\n", testId, queryCollector.Queries.ToList()[i], resultQuery)

			}
		}
	}
}

func testMigrations() []*migration.Migration {
	return []*migration.Migration{
		migration.NewMigration(
			"000000000000",
			func(migration *migration.Migration) error {
				migration.QueryBag.AddQuery("UPDATE table1 SET col1 = \"val\"")
				migration.QueryBag.AddQuery("UPDATE table1 SET col2 = \"val\"")
				return nil
			},
			func(migration *migration.Migration) error {
				migration.QueryBag.AddQuery("UPDATE table1 SET col2 = \"old val\"")
				migration.QueryBag.AddQuery("UPDATE table1 SET col1 = \"old val\"")
				return nil
			}),
		migration.NewMigration(
			"202102101045",
			func(migration *migration.Migration) error {
				migration.QueryBag.AddQuery("UPDATE table2 SET col1 = \"val\"")
				return nil
			},
			func(migration *migration.Migration) error {
				migration.QueryBag.AddQuery("UPDATE table2 SET col1 = \"old val\"")
				return nil
			}),
		migration.NewMigration(
			"202102150011",
			func(migration *migration.Migration) error {
				migration.QueryBag.AddQuery("UPDATE table3 SET col1 = \"val\"")
				return nil
			},
			func(migration *migration.Migration) error {
				migration.QueryBag.AddQuery("UPDATE table3 SET col1 = \"old val\"")
				return nil
			}),
		migration.NewMigration(
			"202103172206",
			func(migration *migration.Migration) error {
				return nil
			},
			func(migration *migration.Migration) error {
				return nil
			}),
	}
}
