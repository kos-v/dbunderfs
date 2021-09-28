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

import "github.com/kos-v/dbunderfs/src/db/migration"

func migration202109101200() *migration.Migration {
	return migration.NewMigration(
		"202109101200",
		func(migration *migration.Migration) error {
			migration.QueryBag.AddQuery(`
				CREATE TABLE {%t_prefix%}descriptors (
					inode      bigint(20) unsigned NOT NULL AUTO_INCREMENT,
					parent     bigint(20) unsigned          DEFAULT NULL,
					name       varchar(255)        NOT NULL,
					type       enum ('DIR','FILE') NOT NULL,
					size       bigint(20)          NOT NULL DEFAULT '0',
					permission varchar(4)          NOT NULL,
					uid        int(11) unsigned    NOT NULL,
					gid        int(11) unsigned    NOT NULL,
					fast_block longblob,
					PRIMARY KEY (inode),
					UNIQUE INDEX ` + "`UNQ__{%t_prefix%}descriptors__parent-name`" + ` (parent, name),
					CONSTRAINT ` + "`FK__{%t_prefix%}descriptors-parent__{%t_prefix%}descriptors-inode`" + ` 
						FOREIGN KEY (parent) REFERENCES {%t_prefix%}descriptors (inode)
							ON UPDATE CASCADE ON DELETE CASCADE
				) ENGINE = InnoDB
				DEFAULT CHARSET = utf8`,
			)

			return nil
		},
		func(migration *migration.Migration) error {
			migration.QueryBag.AddQuery(`DROP TABLE {%t_prefix%}descriptors`)

			return nil
		})
}
