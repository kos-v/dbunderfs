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

import (
	"github.com/kos-v/dbunderfs/src/container"
	"github.com/kos-v/dbunderfs/src/db/migration"
	"github.com/kos-v/dbunderfs/test/helpers/db"
	helpers "github.com/kos-v/dbunderfs/test/helpers/db/migration"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func CreateMigrator(dir string, downNumber int, migrations []*migration.Migration, needUp []string) *migration.Migrator {
	commiterStorage := container.Collection{}
	for _, mig := range migrations {
		for _, upId := range needUp {
			if mig.Id == upId {
				commiterStorage.Append(upId)
			}
		}
	}

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	return &migration.Migrator{
		Commiter:            &helpers.CommiterStub{Storage: &commiterStorage},
		DownDirectionNumber: downNumber,
		Direction:           dir,
		QueryExecutor:       &db.QueryExecuteCollector{Queries: &container.Collection{}},
		Logger:              logger,
		Migrations:          migrations,
	}
}
