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
	"github.com/kos-v/dbunderfs/src/container"
)

const (
	RootName string = "/"
)

type DBInstance interface {
	Close() error
	Connect() (*sql.DB, error)
	GetDriverName() string
	GetPool() *sql.DB
	HasConnection() bool
	Reconnect() (*sql.DB, error)
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

type DBFactory interface {
	CreateDataBlockRepository() DataBlockRepository
	CreateDescriptorRepository() DescriptorRepository
}