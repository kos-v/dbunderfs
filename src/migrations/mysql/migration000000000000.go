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

package mysql

import (
	"github.com/kos-v/dbunderfs/src/db/migration"
)

func migration000000000000() *migration.Migration {
	return migration.NewMigration(
		"000000000000",
		func(migration *migration.Migration) error {
			migration.QueryBag.AddQuery(`
				CREATE TABLE {%t_prefix%}migrations (
					id varchar(255) NOT NULL,
					migrated_at int(11) NOT NULL,
					PRIMARY KEY (id)
				) ENGINE = InnoDB
				DEFAULT CHARSET = utf8`,
			)

			return nil
		},
		func(migration *migration.Migration) error {
			migration.QueryBag.AddQuery(`DROP TABLE {%t_prefix%}migrations`)

			return nil
		})
}