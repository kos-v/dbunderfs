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
	"github.com/kos-v/dbunderfs/internal/container"
)

const (
	RootName string = "/"
)

type Instance interface {
	Close() error
	Connect() (*sql.DB, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	GetDriverName() string
	GetDSN() DSN
	GetPool() *sql.DB
	HasConnection() bool
	Reconnect() (*sql.DB, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type DSN interface {
	GetDatabase() string
	GetTablePrefix() string
	ToString() string
}

type DataBlockRepository interface {
	FindFirst(descr DescriptorInterface) (DataBlockNodeInterface, error)
	Write(descr DescriptorInterface, data *[]byte) error
}

type DescriptorRepository interface {
	Create(parent Inode, name string, dType DescriptorType, attrs DescriptorAttrs) (DescriptorInterface, error)
	FindChildrenByInode(parentInode Inode) (container.CollectionInterface, error)
	FindRoot() (DescriptorInterface, error)
	FindSingleByInode(inode Inode) (DescriptorInterface, error)
	FindSingleByName(parent Inode, target string) (DescriptorInterface, error)
	IsExistsByName(parent Inode, name string) (bool, error)
	RemoveByName(parent Inode, name string) error
}

type RepositoryRegistry interface {
	GetDataBlockRepository() DataBlockRepository
	GetDescriptorRepository() DescriptorRepository
}

type QueryExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}
