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
	log "github.com/sirupsen/logrus"
	"sort"
)

type Migrator struct {
	Commiter            Commiter
	Direction           string
	DownDirectionNumber int
	Instance            db.DBInstance
	Logger              *log.Logger
	Migrations          []*Migration
}

func (m *Migrator) Migrate() error {
	m.prepare()

	switch m.Direction {
	case DirUp:
		return m.up()
	case DirDown:
		return m.down()
	}

	return fmt.Errorf("unexpected value %q for direction param. Should be %q or %q", m.Direction, DirUp, DirDown)
}

func (m *Migrator) down() error {
	migrations := []*Migration{}
	for _, migration := range m.Migrations {
		if ok, err := m.Commiter.IsCommited(migration); err == nil {
			if ok {
				migrations = append(migrations, migration)
			}
		} else {
			return err
		}
	}

	if len(migrations) == 0 {
		m.Logger.Warn("No migrations to rollback.")
		return nil
	}

	downNumber := m.DownDirectionNumber
	if downNumber < 0 || downNumber > len(migrations) {
		downNumber = len(migrations)
	}

	begin := len(migrations) - 1
	end := begin - downNumber
	for i := begin; i > end; i-- {
		migration := migrations[i]

		m.Logger.Infof("Rollback migration %s...", migration.Id)
		queryBag, err := migration.Exec(DirDown)
		if err != nil {
			m.Logger.Errorf("Fail %s.", migration.Id)
			return err
		}

		for _, query := range queryBag.GetQueries() {
			m.Logger.Infof("Executing query: %s", query)
			_, err = m.Instance.Exec(query)
			if err != nil {
				m.Logger.Error("Fail.")
				return err
			}
		}

		err = m.Commiter.Rollback(migration)
		if err != nil {
			return err
		}

		m.Logger.Infof("%s done.", migration.Id)
	}

	return nil
}

func (m *Migrator) up() error {
	executeCount := 0
	for _, migration := range m.Migrations {
		isCommited, err := m.Commiter.IsCommited(migration)
		if err != nil {
			return err
		}
		if isCommited {
			continue
		}

		m.Logger.Infof("Executing migration %s...", migration.Id)
		queryBag, err := migration.Exec(DirUp)
		if err != nil {
			m.Logger.Errorf("Fail %s.", migration.Id)
			return err
		}

		for _, query := range queryBag.GetQueries() {
			m.Logger.Infof("Executing query: %s", query)
			_, err = m.Instance.Exec(query)
			if err != nil {
				m.Logger.Error("Fail.")
				return err
			}
		}
		err = m.Commiter.Commit(migration)
		if err != nil {
			return err
		}

		executeCount++

		m.Logger.Infof("%s done.", migration.Id)
	}

	if executeCount == 0 {
		m.Logger.Warn("No migrations to execute.")
	}

	return nil
}

func (m *Migrator) prepare() {
	sort.Sort(SortMigrationsByDatetime(m.Migrations))

	for i := 0; i < len(m.Migrations); i++ {
		m.Migrations[i].QueryBag = &QueryBag{Queries: container.Collection{List: []interface{}{}}}
	}
}

