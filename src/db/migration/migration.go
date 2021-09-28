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
	"github.com/kos-v/dbunderfs/src/container"
	"github.com/kos-v/dbunderfs/src/db"
	"time"
)

const DirDown = "down"
const DirUp = "up"

type DirectionHandler func(migration *Migration) error

func NewMigration(id string, onUp DirectionHandler, onDown DirectionHandler) *Migration {
	m := Migration{Id: id}
	m.OnUp(onUp)
	m.OnDown(onDown)

	return &m
}

type Migration struct {
	Id        string
	CreatedAt int64
	Query     Query
	QueryBag  *QueryBag

	onDownHandler DirectionHandler
	onUpHandler   DirectionHandler
}

func (m *Migration) Exec(direction string) (*QueryBag, error) {
	var err error

	switch direction {
	case DirUp:
		if m.onUpHandler != nil {
			err = m.onUpHandler(m)
		}
	case DirDown:
		if m.onDownHandler != nil {
			err = m.onDownHandler(m)
		}
	default:
		return nil, fmt.Errorf("unexpected value %q for direction param. Should be %q or %q", DirUp, DirDown, direction)
	}

	if err != nil {
		return nil, err
	}

	return m.QueryBag, nil
}

func (m *Migration) OnDown(handler DirectionHandler) *Migration {
	m.onDownHandler = handler
	return m
}

func (m *Migration) OnUp(handler DirectionHandler) *Migration {
	m.onUpHandler = handler
	return m
}

type Query struct {
	Instance db.DBInstance
}

type QueryBag struct {
	Queries container.Collection
}

func (qb *QueryBag) AddQuery(query string) {
	qb.Queries.Append(query)
}

func (qb *QueryBag) GetQueries() []string {
	queries := make([]string, qb.Queries.Len())
	for i, query := range qb.Queries.ToList() {
		queries[i] = query.(string)
	}

	return queries
}

type SortMigrationsByDatetime []*Migration

func (m SortMigrationsByDatetime) Len() int           { return len(m) }
func (m SortMigrationsByDatetime) Less(i, j int) bool { return m[i].Id < m[j].Id }
func (m SortMigrationsByDatetime) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

type Commiter interface {
	Commit(migration *Migration) error
	IsCommited(migration *Migration) (bool, error)
	Rollback(migration *Migration) error
}

type MySQLCommiter struct {
	Instance db.DBInstance
}

func (c *MySQLCommiter) Commit(migration *Migration) error {
	migration.CreatedAt = time.Now().Unix()

	query := `INSERT INTO {%t_prefix%}migrations (id, migrated_at) VALUES (?, ?)`
	_, err := c.Instance.Exec(query, migration.Id, migration.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (c *MySQLCommiter) IsCommited(migration *Migration) (bool, error) {
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

func (c *MySQLCommiter) Rollback(migration *Migration) error {
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

func (c *MySQLCommiter) isExistsCommitTable() (bool, error) {
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
