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
