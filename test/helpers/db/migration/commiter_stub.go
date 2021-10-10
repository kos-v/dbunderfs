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
)

type CommiterStub struct {
	Storage *container.Collection
}

func (c *CommiterStub) Commit(migration *migration.Migration) error {
	c.Storage.Append(migration.Id)

	return nil
}

func (c *CommiterStub) IsCommited(migration *migration.Migration) (bool, error) {
	for _, id := range c.Storage.ToList() {
		if migration.Id == id {
			return true, nil
		}
	}

	return false, nil
}

func (c *CommiterStub) Rollback(migration *migration.Migration) error {
	for i, id := range c.Storage.ToList() {
		if id == migration.Id {
			c.Storage.Remove(i)
			break
		}
	}

	return nil
}
