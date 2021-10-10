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
	"fmt"
	"github.com/kos-v/dbunderfs/internal/container"
	"github.com/kos-v/dbunderfs/internal/db/migration"
	"testing"
)

func TestMigration_Up(t *testing.T) {
	testMigrationByDirection(t, migration.DirUp)
}

func TestMigration_Down(t *testing.T) {
	testMigrationByDirection(t, migration.DirDown)
}

func TestMigration_IncorrectDirection(t *testing.T) {
	testDirection := "foo"
	expectedError := fmt.Errorf("unexpected value %q for direction param. Should be %q or %q", migration.DirUp, migration.DirDown, "foo")

	m := migration.NewMigration(
		"000000000000",
		func(m *migration.Migration) error {
			m.QueryBag.AddQuery("SELECT 1")
			return nil
		},
		func(m *migration.Migration) error {
			m.QueryBag.AddQuery("SELECT 1")
			return nil
		})

	_, err := m.Exec(testDirection)
	if err == nil {
		t.Fatalf("The method did not return error")
	}
	if err.Error() != expectedError.Error() {
		t.Fatalf("The method did not return the error that was expected.\nExpected: %q.\nResult: %q", expectedError.Error(), err.Error())
	}
}

func testMigrationByDirection(t *testing.T, dir string) {
	tests := []struct {
		queries  []string
		expected []string
	}{
		{queries: []string{}, expected: []string{}},
		{queries: []string{"SELECT 1", "SELECT 2"}, expected: []string{"SELECT 1", "SELECT 2"}},
		{queries: []string{"SELECT 1", "SELECT 2", "SELECT 3"}, expected: []string{"SELECT 1", "SELECT 2", "SELECT 3"}},
	}

	for testId, test := range tests {
		m := migration.NewMigration(
			"000000000000",
			func(m *migration.Migration) error {
				for _, query := range test.queries {
					m.QueryBag.AddQuery(query)
				}
				return nil
			},
			func(m *migration.Migration) error {
				for _, query := range test.queries {
					m.QueryBag.AddQuery(query)
				}
				return nil
			})
		m.QueryBag = &migration.QueryBag{Queries: container.Collection{}}

		result, resultErr := m.Exec(dir)
		if resultErr != nil {
			t.Fatalf("Test %v fail: the method returned an unexpected error. Error: %e", testId, resultErr)
		}

		if len(result.GetQueries()) != len(test.expected) {
			t.Fatalf("Test %v fail: the number of items is not as expected.\nExpected: %v. Result: %v.\n", testId, len(test.expected), len(result.GetQueries()))
		}

		for i, resultQuery := range result.GetQueries() {
			if resultQuery != test.expected[i] {
				t.Fatalf("Test %v fail: result data is not as expected.\nExpected: %v. Result: %v.\n", testId, test.expected[i], resultQuery)
			}
		}
	}
}
