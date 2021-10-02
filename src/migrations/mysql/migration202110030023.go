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

func migration202110030023() *migration.Migration {
	return migration.NewMigration(
		"202110030023",
		func(migration *migration.Migration) error {
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

			return nil
		},
		func(migration *migration.Migration) error {
			migration.QueryBag.AddQuery(`DROP PROCEDURE findDescriptorByPath`)

			return nil
		})
}
