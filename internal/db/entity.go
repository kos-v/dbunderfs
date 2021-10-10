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

package db

import (
	"database/sql"
	"io/fs"
	"strconv"
	"sync"
)

type Inode uint64

type DescriptorType string

const (
	DT_Dir  DescriptorType = "DIR"
	DT_File DescriptorType = "FILE"
)

type DataBlockNodeInterface interface {
	Add(offset uint64, newData *[]byte) int
	GetData() *[]byte
}

type DataBlockNode struct {
	mu   sync.Mutex
	Data []byte
}

func (d *DataBlockNode) Add(offset uint64, newData *[]byte) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	if int(offset) > len(d.Data) {
		offset = uint64(len(d.Data))
	}

	newLen := len(d.Data) + len(*newData) - (len(d.Data) - int(offset))
	data := make([]byte, newLen)
	copy(data, d.Data)

	addedLen := copy(data[offset:], *newData)
	d.Data = data

	return addedLen
}

func (d *DataBlockNode) GetData() *[]byte {
	d.mu.Lock()
	defer d.mu.Unlock()

	return &d.Data
}

type DescriptorInterface interface {
	GetInode() Inode
	GetName() string
	GetParent() Inode
	GetPermission() fs.FileMode
	GetSize() uint64
	GetType() DescriptorType
	GetUID() uint32
	GetGID() uint32
	IsRoot() bool
}

type DescriptorAttrs struct {
	Size       uint64
	Permission string
	UID        uint32
	GID        uint32
}

type Descriptor struct {
	DescriptorAttrs

	Inode  Inode
	Parent sql.NullInt64
	Name   string
	Type   DescriptorType
}

func (d *Descriptor) GetInode() Inode {
	return d.Inode
}

func (d *Descriptor) GetName() string {
	return d.Name
}

func (d *Descriptor) GetParent() Inode {
	return Inode(d.Parent.Int64)
}

func (d *Descriptor) GetPermission() fs.FileMode {
	octal, err := strconv.ParseUint(d.Permission, 8, 16)
	if err != nil {
		return 0
	}

	return fs.FileMode(octal)
}

func (d *Descriptor) GetSize() uint64 {
	return d.Size
}

func (d *Descriptor) GetType() DescriptorType {
	return d.Type
}

func (d *Descriptor) GetUID() uint32 {
	return d.UID
}

func (d *Descriptor) GetGID() uint32 {
	return d.GID
}

func (d *Descriptor) IsRoot() bool {
	val, _ := d.Parent.Value()
	return val == nil
}
