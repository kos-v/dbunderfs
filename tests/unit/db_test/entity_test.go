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

import "testing"
import "github.com/kos-v/dbunderfs/src/db"

func TestDataBlockNode_Add(t *testing.T) {
	tests := []struct {
		defaultData    []byte
		offset         uint64
		addData        []byte
		expectedData   []byte
		expectedAddLen int
	}{
		{[]byte{}, 0, []byte{}, []byte{}, 0},
		{[]byte{}, 0, []byte{1, 2, 3}, []byte{1, 2, 3}, 3},
		{[]byte{}, 1, []byte{1, 2, 3}, []byte{1, 2, 3}, 3},
		{[]byte{1, 2, 3}, 0, []byte{4, 5, 6}, []byte{4, 5, 6}, 3},
		{[]byte{1, 2, 3}, 1, []byte{4, 5, 6}, []byte{1, 4, 5, 6}, 3},
		{[]byte{1, 2, 3}, 2, []byte{4, 5, 6}, []byte{1, 2, 4, 5, 6}, 3},
		{[]byte{1, 2, 3}, 3, []byte{4, 5, 6}, []byte{1, 2, 3, 4, 5, 6}, 3},
		{[]byte{1, 2, 3}, 4, []byte{4, 5, 6}, []byte{1, 2, 3, 4, 5, 6}, 3},
	}

	for id, test := range tests {
		id += 1
		dataBlock := db.DataBlockNode{Data: test.defaultData}
		addLen := dataBlock.Add(test.offset, &test.addData)
		if addLen != test.expectedAddLen {
			t.Errorf("Test %v fail: the number of items added is not as expected.\nExpected: %v. Result: %v.\n", id, test.expectedAddLen, addLen)
		}
		if len(*dataBlock.GetData()) != len(test.expectedData) {
			t.Errorf("Test %v fail: object contains unexpected number of items.\nExpected: %v. Result: %v.\n", id, len(test.expectedData), len(*dataBlock.GetData()))
		}
		for i, item := range *dataBlock.GetData() {
			if item != test.expectedData[i] {
				t.Errorf("Test %v fail: result data is not as expected.\nExpected: %v. Result: %v.\n", id, test.expectedData, *dataBlock.GetData())
			}
		}
	}
}
