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
	"github.com/kos-v/dbunderfs/src/db"
	"github.com/kos-v/dbunderfs/src/db/migration"
	"time"
)

type Commiter struct {
	Instance db.Instance
}

func (c *Commiter) Commit(migration *migration.Migration) error {
	migration.CreatedAt = time.Now().Unix()

	query := `INSERT INTO {%t_prefix%}migrations (id, migrated_at) VALUES (?, ?)`
	_, err := c.Instance.Exec(query, migration.Id, migration.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commiter) IsCommited(migration *migration.Migration) (bool, error) {
	if ok, err := c.isExistsCommitTable(); err == nil {
		if !ok {
			return false, nil
		}
	} else {
		return false, err
	}

	var isExistsMigration int
	row := c.Instance.QueryRow("SELECT COUNT(*) FROM {%t_prefix%}migrations WHERE id = ?", migration.Id)
	err := row.Scan(&isExistsMigration)
	if err != nil {
		return false, err
	}

	return isExistsMigration > 0, nil
}

func (c *Commiter) Rollback(migration *migration.Migration) error {
	if ok, err := c.isExistsCommitTable(); err == nil {
		if !ok {
			return nil
		}
	} else {
		return err
	}

	_, err := c.Instance.Exec("DELETE FROM {%t_prefix%}migrations WHERE id = ?", migration.Id)
	return err
}

func (c *Commiter) isExistsCommitTable() (bool, error) {
	var isExists int

	query := `SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_schema = ? AND table_name = '{%t_prefix%}migrations' 
		LIMIT 1`
	row := c.Instance.QueryRow(query, c.Instance.GetDSN().GetDatabase())

	err := row.Scan(&isExists)
	if err != nil {
		return false, err
	}

	return isExists != 0, nil
}
