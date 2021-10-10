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

package migration

import (
	"github.com/kos-v/dbunderfs/internal/db/migration"
	"testing"
)

func TestQueryBag(t *testing.T) {
	tests := []struct {
		queries  []string
		expected []string
	}{
		{queries: []string{}, expected: []string{}},
		{queries: []string{"SELECT 1"}, expected: []string{"SELECT 1"}},
		{queries: []string{"SELECT 1", "SELECT 2", "SELECT 3"}, expected: []string{"SELECT 1", "SELECT 2", "SELECT 3"}},
	}

	for testId, test := range tests {
		qb := migration.QueryBag{}
		for _, query := range test.queries {
			qb.AddQuery(query)
		}

		if len(qb.GetQueries()) != len(test.expected) {
			t.Fatalf("Test %v fail: the number of items is not as expected.\nExpected: %v. Result: %v.\n", testId, len(test.expected), len(qb.GetQueries()))
		}
	}
}
