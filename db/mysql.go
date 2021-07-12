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
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kos-v/dbunderfs/container"
	"sync"
)

type MySQLInstance struct {
	DSN string

	sync.RWMutex
	pool *sql.DB
}

func (inst *MySQLInstance) Close() error {
	if inst.pool == nil {
		return nil
	}
	return inst.pool.Close()
}

func (inst *MySQLInstance) Connect() (*sql.DB, error) {
	if inst.HasConnection() {
		return inst.pool, nil
	}

	pool, err := inst.generatePool(inst.DSN)
	if err != nil {
		return nil, err
	}

	inst.Lock()
	inst.pool = pool
	inst.Unlock()

	return inst.pool, nil
}

func (inst *MySQLInstance) GetDriverName() string {
	return "mysql"
}

func (inst *MySQLInstance) GetPool() *sql.DB {
	return inst.pool
}

func (inst *MySQLInstance) HasConnection() bool {
	if inst.pool != nil {
		if err := inst.pool.Ping(); err == nil {
			return true
		}
	}

	return false
}

func (inst *MySQLInstance) Reconnect() (*sql.DB, error) {
	pool, err := inst.generatePool(inst.DSN)
	if err != nil {
		return nil, err
	}

	inst.Lock()
	defer inst.Unlock()

	if inst.HasConnection() {
		if err = inst.Close(); err != nil {
			return nil, err
		}
	}
	inst.pool = pool

	return inst.pool, nil
}

func (inst *MySQLInstance) generatePool(dsn string) (*sql.DB, error) {
	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(); err != nil {
		return nil, err
	}

	return pool, err
}

type MySQLDataBlockRepository struct {
	instance DBInstance
}

func (repo *MySQLDataBlockRepository) FindFirst(descr DescriptorInterface) (DataBlockNodeInterface, error) {
	row := repo.instance.GetPool().QueryRow(`SELECT fast_block FROM descriptors WHERE inode = ?`, descr.GetInode())

	dataBlock := DataBlockNode{data: &[]byte{}}
	err := row.Scan(dataBlock.data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &dataBlock, nil
}

type MySQLDescriptorRepository struct {
	instance DBInstance
}

func (dr *MySQLDescriptorRepository) Create(parent Inode, name string, dType DescriptorType, attrs DescriptorAttrs) (DescriptorInterface, error) {
	sqlStatement := `
	INSERT INTO descriptors (parent, name, type, size, permission,  uid, gid)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := dr.instance.GetPool().Exec(sqlStatement,
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

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	descr, err := dr.FindSingleByInode(Inode(lastInsertId))
	if err != nil {
		return nil, err
	}

	return descr, nil
}

func (dr *MySQLDescriptorRepository) FindChildrenByInode(parentInode Inode) (container.CollectionInterface, error) {
	rows, err := dr.instance.GetPool().Query(`
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

func (dr *MySQLDescriptorRepository) FindRoot() (DescriptorInterface, error) {
	row := dr.instance.GetPool().QueryRow(`CALL findDescriptorByPath(?, NULL, 1)`, RootName)

	descr, err := dr.hydrateDescriptor(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if descr.GetType() != DT_Dir {
		return nil, nil
	}

	return descr, nil
}

func (dr *MySQLDescriptorRepository) FindSingleByInode(inode Inode) (DescriptorInterface, error) {
	row := dr.instance.GetPool().QueryRow(`
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

func (dr *MySQLDescriptorRepository) FindSingleByName(parent Inode, target string) (DescriptorInterface, error) {
	row := dr.instance.GetPool().QueryRow(`
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

func (dr *MySQLDescriptorRepository) IsExistsByName(parent Inode, name string) (bool, error) {
	descr, err := dr.FindSingleByName(parent, name)
	if err != nil {
		return false, err
	}

	return descr != nil, nil
}

func (dr *MySQLDescriptorRepository) hydrateDescriptor(row interface{}) (DescriptorInterface, error) {
	descr := Descriptor{}

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

type MySQLFactory struct {
	Instance DBInstance
}

func (f *MySQLFactory) CreateDataBlockRepository() DataBlockRepository {
	return &MySQLDataBlockRepository{instance: f.Instance}
}

func (f *MySQLFactory) CreateDescriptorRepository() DescriptorRepository {
	return &MySQLDescriptorRepository{instance: f.Instance}
}
