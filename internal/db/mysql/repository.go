/*
   Copyright 2021 The DbunderFS Contributors.

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
	"database/sql"
	"errors"
	"fmt"
	"github.com/kos-v/dbunderfs/internal/container"
	"github.com/kos-v/dbunderfs/internal/db"
)

type DataBlockRepository struct {
	instance db.Instance
}

func (repo *DataBlockRepository) FindFirst(descr db.DescriptorInterface) (db.DataBlockNodeInterface, error) {
	row := repo.instance.QueryRow(`SELECT fast_block FROM descriptors WHERE inode = ?`, descr.GetInode())

	dataBlock := db.DataBlockNode{Data: []byte{}}
	err := row.Scan(&dataBlock.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &dataBlock, nil
}

func (repo *DataBlockRepository) Write(descr db.DescriptorInterface, data *[]byte) error {
	_, err := repo.instance.Exec(`UPDATE descriptors SET fast_block = ?, size = ? WHERE inode = ?`,
		*data,
		len(*data),
		descr.GetInode(),
	)

	return err
}

type DescriptorRepository struct {
	instance db.Instance
}

func (dr *DescriptorRepository) Create(parent db.Inode, name string, dType db.DescriptorType, attrs db.DescriptorAttrs) (db.DescriptorInterface, error) {
	sqlStatement := `
	INSERT INTO descriptors (parent, name, type, size, permission,  uid, gid)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := dr.instance.Exec(sqlStatement,
		parent,
		name,
		dType,
		attrs.Size,
		attrs.Permission,
		attrs.UID,
		attrs.GID,
	)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	descr, err := dr.FindSingleByInode(db.Inode(lastInsertID))
	if err != nil {
		return nil, err
	}

	return descr, nil
}

func (dr *DescriptorRepository) FindChildrenByInode(parentInode db.Inode) (container.CollectionInterface, error) {
	rows, err := dr.instance.Query(`
		SELECT inode, parent, name, type, size, permission,  uid, gid 
		FROM descriptors
		WHERE parent = ?
		ORDER BY type, name`, parentInode,
	)
	if err != nil {
		return nil, err
	}

	collection := container.Collection{}
	for rows.Next() {
		descr, err := dr.hydrateDescriptor(rows)
		if err != nil {
			return nil, err
		}

		collection.Append(descr)
	}

	return &collection, nil
}

func (dr *DescriptorRepository) FindRoot() (db.DescriptorInterface, error) {
	row := dr.instance.QueryRow(`CALL findDescriptorByPath(?, NULL, 1)`, db.RootName)

	descr, err := dr.hydrateDescriptor(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if descr.GetType() != db.DT_Dir {
		return nil, nil
	}

	return descr, nil
}

func (dr *DescriptorRepository) FindSingleByInode(inode db.Inode) (db.DescriptorInterface, error) {
	row := dr.instance.QueryRow(`
		SELECT inode, parent, name, type, size, permission,  uid, gid 
		FROM descriptors 
		WHERE inode = ?`, inode,
	)

	descr, err := dr.hydrateDescriptor(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return descr, nil
}

func (dr *DescriptorRepository) FindSingleByName(parent db.Inode, target string) (db.DescriptorInterface, error) {
	row := dr.instance.QueryRow(`
		SELECT inode, parent, name, type, size, permission,  uid, gid
		FROM descriptors 
		WHERE parent = ?  AND name = ?`, parent, target,
	)

	descr, err := dr.hydrateDescriptor(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return descr, nil
}

func (dr *DescriptorRepository) IsExistsByName(parent db.Inode, name string) (bool, error) {
	descr, err := dr.FindSingleByName(parent, name)
	if err != nil {
		return false, err
	}

	return descr != nil, nil
}

func (dr *DescriptorRepository) RemoveByName(parent db.Inode, name string) error {
	exists, err := dr.IsExistsByName(parent, name)
	if err != nil {
		return err
	}

	if exists == false {
		return fmt.Errorf("Node %s was not found in parent %d", name, parent)
	}

	_, err = dr.instance.Exec("DELETE FROM descriptors WHERE parent = ?  AND name = ?", parent, name)
	return err
}

func (dr *DescriptorRepository) hydrateDescriptor(row interface{}) (db.DescriptorInterface, error) {
	descr := db.Descriptor{}

	var err error
	switch row.(type) {
	case *sql.Row:
		row := row.(*sql.Row)
		err = row.Scan(
			&descr.Inode,
			&descr.Parent,
			&descr.Name,
			&descr.Type,
			&descr.Size,
			&descr.Permission,
			&descr.UID,
			&descr.GID,
		)
	case *sql.Rows:
		row := row.(*sql.Rows)
		err = row.Scan(
			&descr.Inode,
			&descr.Parent,
			&descr.Name,
			&descr.Type,
			&descr.Size,
			&descr.Permission,
			&descr.UID,
			&descr.GID,
		)
	default:
		return nil, fmt.Errorf("type \"%T\" not support for row argument", row)
	}

	if err != nil {
		return nil, err
	}

	return &descr, nil
}
