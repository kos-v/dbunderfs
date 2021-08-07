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

/*!40101 SET @OLD_CHARACTER_SET_CLIENT = @@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS = 0 */;
/*!40101 SET @OLD_SQL_MODE = @@SQL_MODE, SQL_MODE = 'NO_AUTO_VALUE_ON_ZERO' */;

CREATE TABLE IF NOT EXISTS `descriptors`
(
    `inode`      bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `parent`     bigint(20) unsigned          DEFAULT NULL,
    `name`       varchar(255)        NOT NULL,
    `type`       enum ('DIR','FILE') NOT NULL,
    `size`       bigint(20)          NOT NULL DEFAULT '0',
    `permission` varchar(4)          NOT NULL,
    `uid`        int(11) unsigned    NOT NULL,
    `gid`        int(11) unsigned    NOT NULL,
    `fast_block` longblob,
    PRIMARY KEY (`inode`),
    UNIQUE KEY `parent-name_uniq` (`parent`, `name`),
    CONSTRAINT `FK--descriptors-parent--descriptors-inode` FOREIGN KEY (`parent`) REFERENCES `descriptors` (`inode`)
        ON UPDATE CASCADE ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

DELIMITER //
CREATE
    DEFINER = `root`@`%` PROCEDURE `findDescriptorByPath`(
    IN `path` VARCHAR(255),
    IN `parent` INT,
    IN `callIndex` INT
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
        SELECT `inode` INTO @subpathId FROM descriptors WHERE `parent` IS NULL AND `name` = @subpath LIMIT 1;
    ELSE
        SELECT `inode` INTO @subpathId FROM descriptors WHERE `parent` = parent AND `name` = @subpath LIMIT 1;
    END IF;

    IF @subpathId IS NULL THEN
        LEAVE this_proc;
    END IF;

    IF callIndex >= @maxDepth THEN
        IF @pathIsRoot = TRUE THEN
            SELECT `inode`,
                   `parent`,
                   `name`,
                   `type`,
                   `size`,
                   `permission`,
                   `uid`,
                   `gid`
            FROM descriptors
            WHERE `parent` IS NULL
              AND `name` = @subpath
            LIMIT 1;
        ELSE
            SELECT `inode`,
                   `parent`,
                   `name`,
                   `type`,
                   `size`,
                   `permission`,
                   `uid`,
                   `gid`
            FROM descriptors
            WHERE `parent` = parent
              AND `name` = @subpath
            LIMIT 1;
        END IF;
    ELSE
        CALL findDescriptorByPath(path, @subpathId, callIndex + 1);
    END IF;
END//
DELIMITER ;

/*!40101 SET SQL_MODE = IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS = IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT = @OLD_CHARACTER_SET_CLIENT */;
