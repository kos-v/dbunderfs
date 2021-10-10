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
	"github.com/kos-v/dbunderfs/src/db/migration"
	"sort"
	"testing"
)

func TestSortMigrationsByDatetime(t *testing.T) {
	createMigrations := func(ids []string) []*migration.Migration {
		migrations := make([]*migration.Migration, len(ids))
		for i, id := range ids {
			migrations[i] = migration.NewMigration(id, nil, nil)
		}

		return migrations
	}

	tests := []struct {
		migrationIds []string
		expected     []string
	}{
		{migrationIds: []string{}, expected: []string{}},
		{
			migrationIds: []string{"000000000000"},
			expected:     []string{"000000000000"},
		},
		{
			migrationIds: []string{"000000000001", "100000000000", "000000000000"},
			expected:     []string{"000000000000", "000000000001", "100000000000"},
		},
		{
			migrationIds: []string{"202102101045", "202201062211", "000007160255", "000000000000"},
			expected:     []string{"000000000000", "000007160255", "202102101045", "202201062211"},
		},
	}

	for testId, test := range tests {
		sortedMigrations := createMigrations(test.migrationIds)
		sort.Sort(migration.SortMigrationsByDatetime(sortedMigrations))

		if len(sortedMigrations) != len(test.expected) {
			t.Fatalf("Test %v fail: the number of items is not as expected.\nExpected: %v. Result: %v.\n", testId, len(test.expected), len(sortedMigrations))
		}

		for i, sortedMigration := range sortedMigrations {
			if sortedMigration.Id != test.expected[i] {
				t.Fatalf("Test %v fail: result data is not as expected.\nExpected: %v. Result: %v.\n", testId, test.expected, sortedMigrations)
			}
		}
	}
}
