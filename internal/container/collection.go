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

package container

import "sync"

type CollectionInterface interface {
	Append(item interface{})
	Len() int
	ToList() []interface{}
}

type Collection struct {
	mu sync.Mutex

	List []interface{}
}

func (coll *Collection) Append(item interface{}) {
	coll.mu.Lock()
	coll.List = append(coll.List, item)
	coll.mu.Unlock()
}

func (coll *Collection) Len() int {
	coll.mu.Lock()
	defer coll.mu.Unlock()

	return len(coll.List)
}

func (coll *Collection) Remove(index int) {
	coll.mu.Lock()
	defer coll.mu.Unlock()

	if index < 0 || index >= len(coll.List) {
		return
	}

	coll.List = append(coll.List[:index], coll.List[index+1:]...)
}

func (coll *Collection) ToList() []interface{} {
	coll.mu.Lock()
	defer coll.mu.Unlock()

	return coll.List
}
