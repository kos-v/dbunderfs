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
	"fmt"
	"github.com/kos-v/dbunderfs/src/db"
	"github.com/kos-v/dbunderfs/src/db/migration"
	"os/user"
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

			migration.QueryBag.AddQuery(`
				CREATE PROCEDURE findDescriptorByPath (
					IN path VARCHAR(255),
					IN parent INT,
					IN callIndex INT
				)
					READS SQL DATA
				this_proc:
				BEGIN
					SET max_sp_recursion_depth := 2048;
				
					IF path = "" THEN
						LEAVE this_proc;
					END IF;
				
					SET @pathIsRoot := parent IS NULL;
				
					SET @maxDepth := 1;
					IF path <> "/" THEN
						SET @maxDepth := ROUND((CHAR_LENGTH(path) - CHAR_LENGTH(REPLACE(path, '/', ""))) / CHAR_LENGTH('/')) + 1;
					END IF;
				
					SET @subpath = REPLACE(SUBSTRING(SUBSTRING_INDEX(path, '/', callIndex),
													 CHAR_LENGTH(SUBSTRING_INDEX(path, '/', callIndex - 1)) + 1), '/', '');
					IF @subpath = "" THEN
						IF @pathIsRoot = TRUE THEN
							SET @subpath := "/";
						ELSE
							LEAVE this_proc;
						END IF;
					END IF;
				
					SET @subpathId := NULL;
					IF @pathIsRoot = TRUE THEN
						SELECT inode INTO @subpathId FROM descriptors WHERE parent IS NULL AND name = @subpath LIMIT 1;
					ELSE
						SELECT inode INTO @subpathId FROM descriptors WHERE parent = parent AND name = @subpath LIMIT 1;
					END IF;
				
					IF @subpathId IS NULL THEN
						LEAVE this_proc;
					END IF;
				
					IF callIndex >= @maxDepth THEN
						IF @pathIsRoot = TRUE THEN
							SELECT inode,
								   parent,
								   name,
								   type,
								   size,
								   permission,
								   uid,
								   gid
							FROM descriptors
							WHERE parent IS NULL
							  AND name = @subpath
							LIMIT 1;
						ELSE
							SELECT inode,
								   parent,
								   name,
								   type,
								   size,
								   permission,
								   uid,
								   gid
							FROM descriptors
							WHERE parent = parent
							  AND name = @subpath
							LIMIT 1;
						END IF;
					ELSE
						CALL findDescriptorByPath(path, @subpathId, callIndex + 1);
					END IF;
				END`,
			)

			currentUser, err := user.Current()
			if err != nil {
				return err
			}

			migration.QueryBag.AddQuery(fmt.Sprintf(
				`INSERT INTO {%%t_prefix%%}descriptors (parent, name, type, size, permission,  uid, gid) VALUES (NULL, %q, %q, 0, 755, %s, %s)`,
				db.RootName,
				string(db.DT_Dir),
				currentUser.Uid,
				currentUser.Gid,
			))

			return nil
		},
		func(migration *migration.Migration) error {
			migration.QueryBag.AddQuery(`DROP PROCEDURE findDescriptorByPath`)
			migration.QueryBag.AddQuery(`DROP TABLE {%t_prefix%}descriptors`)
			migration.QueryBag.AddQuery(`DROP TABLE {%t_prefix%}migrations`)

			return nil
		})
}
