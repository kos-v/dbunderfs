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

package container

import "testing"
import "github.com/kos-v/dbunderfs/src/container"

func TestCollection_Append(t *testing.T) {
	toInterfaceList := func(items []string) []interface{} {
		interfaceItems := make([]interface{}, len(items))
		for i, item := range items {
			interfaceItems[i] = item
		}
		return interfaceItems
	}

	tests := []struct {
		initItems     []string
		appendItems   []string
		expectedItems []string
	}{
		{[]string{}, []string{"foo", "bar", "baz"}, []string{"foo", "bar", "baz"}},
		{[]string{"foo"}, []string{"bar", "baz"}, []string{"foo", "bar", "baz"}},
		{[]string{"foo", "bar"}, []string{"baz"}, []string{"foo", "bar", "baz"}},
		{[]string{"foo", "bar", "baz"}, []string{}, []string{"foo", "bar", "baz"}},
	}

	for id, test := range tests {
		coll := container.Collection{List: toInterfaceList(test.initItems)}
		for _, appendItem := range test.appendItems {
			coll.Append(appendItem)
		}

		if coll.Len() != len(test.expectedItems) {
			t.Errorf("Test %v fail: object contains unexpected number of items.\nExpected: %v. Result: %v.\n", id, len(test.expectedItems), coll.Len())
			return
		}

		for i, item := range coll.ToList() {
			if item != test.expectedItems[i] {
				t.Errorf("Test %v fail: result data is not as expected.\nExpected: %v. Result: %v.\n", id, test.expectedItems[i], item.(string))
			}
		}
	}
}

func TestCollection_Remove(t *testing.T) {
	tests := []struct {
		initItems     []string
		removeIndex   int
		expectedItems []string
	}{
		{[]string{}, 0, []string{}},
		{[]string{}, 100, []string{}},
		{[]string{"foo"}, 0, []string{}},
		{[]string{"foo", "bar"}, 1, []string{"foo"}},
		{[]string{"foo", "bar", "baz"}, 1, []string{"foo", "baz"}},
		{[]string{"foo", "bar", "baz"}, 0, []string{"bar", "baz"}},
		{[]string{"foo"}, 1, []string{"foo"}},
	}

	for id, test := range tests {
		coll := container.Collection{List: []interface{}{}}
		for _, appendItem := range test.initItems {
			coll.Append(appendItem)
		}

		coll.Remove(test.removeIndex)

		if coll.Len() != len(test.expectedItems) {
			t.Errorf("Test %v fail: object contains unexpected number of items.\nExpected: %v. Result: %v.\n", id, len(test.expectedItems), coll.Len())
			return
		}

		for i, item := range coll.ToList() {
			if item != test.expectedItems[i] {
				t.Errorf("Test %v fail: result data is not as expected.\nExpected: %v. Result: %v.\n", id, test.expectedItems[i], item.(string))
			}
		}
	}
}
