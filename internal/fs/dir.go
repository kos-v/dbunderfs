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

package fs

import (
	"bazil.org/fuse"
	fuseFS "bazil.org/fuse/fs"
	"fmt"
	"github.com/kos-v/dbunderfs/internal/db"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
)

type Dir struct {
	descriptor db.DescriptorInterface
	fs         *FS
}

var _ fuseFS.Node = (*Dir)(nil)

func (d *Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	log.Infof("Reads dir attrs of %d:%s", d.descriptor.GetInode(), db.RootName)

	attr.Inode = uint64(d.descriptor.GetInode())
	attr.Mode = os.ModeDir | d.descriptor.GetPermission()
	attr.Uid = d.descriptor.GetUID()
	attr.Gid = d.descriptor.GetGID()
	attr.Size = d.descriptor.GetSize()

	return nil
}

var _ = fuseFS.NodeRequestLookuper(&Dir{})

func (d *Dir) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (fuseFS.Node, error) {
	requestName := req.Name
	log.Infof("Lookup in \"%d:%s\". Request: %s", d.descriptor.GetInode(), d.descriptor.GetName(), requestName)

	repo := d.fs.RepositoryRegistry.GetDescriptorRepository()
	descr, err := repo.FindSingleByName(d.descriptor.GetInode(), requestName)
	if err != nil {
		log.Warnf("Error: %s", err.Error())
		return nil, err
	}

	if descr == nil {
		log.Warnf("Request path not exists")
		return nil, fuse.ENOENT
	}

	if descr.GetType() == db.DT_Dir {
		return &Dir{
			descriptor: descr,
			fs:         d.fs,
		}, nil
	}

	return &File{
		descriptor: descr,
		fs:         d.fs,
	}, nil
}

var _ = fuseFS.HandleReadDirAller(&Dir{})

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Infof("Read dir all in %d:%s", d.descriptor.GetInode(), d.descriptor.GetName())

	collection, err := d.fs.RepositoryRegistry.GetDescriptorRepository().FindChildrenByInode(d.descriptor.GetInode())
	if err != nil {
		return nil, err
	}

	var res []fuse.Dirent
	for _, item := range collection.ToList() {
		item := item.(db.DescriptorInterface)

		var de fuse.Dirent
		if item.GetType() == db.DT_Dir {
			de.Type = fuse.DT_Dir
		} else {
			de.Type = fuse.DT_File
		}
		de.Name = item.GetName()
		res = append(res, de)
	}

	return res, nil
}

var _ fuseFS.NodeMkdirer = (*Dir)(nil)

func (d *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fuseFS.Node, error) {
	log.Infof("Mkdir for %d:%s", d.descriptor.GetInode(), req.Name)

	repo := d.fs.RepositoryRegistry.GetDescriptorRepository()
	isExists, err := repo.IsExistsByName(d.descriptor.GetInode(), req.Name)
	if err != nil {
		log.Warnf("Error: %s: ", err.Error())
		return nil, err
	}
	if isExists {
		msg := fmt.Sprintf("Directory with name \"%s\" already exists", req.Name)
		log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	newDescr, err := repo.Create(d.descriptor.GetInode(), req.Name, db.DT_Dir, db.DescriptorAttrs{
		GID:        req.Gid,
		UID:        req.Uid,
		Permission: "0755",
		Size:       0,
	})
	if err != nil {
		log.Errorf("Error creating directory. Error: %s", err.Error())
		return nil, err
	}

	return &Dir{
		descriptor: newDescr,
		fs:         d.fs,
	}, nil
}

var _ = fuseFS.NodeCreater(&Dir{})

func (d *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fuseFS.Node, fuseFS.Handle, error) {
	log.Infof("Create file %s in %s[%d]", req.Name, d.descriptor.GetName(), d.descriptor.GetInode())

	repo := d.fs.RepositoryRegistry.GetDescriptorRepository()
	isExists, err := repo.IsExistsByName(d.descriptor.GetInode(), req.Name)
	if err != nil {
		log.Warnf("Error: %s: ", err.Error())
		return nil, nil, err
	}
	if isExists {
		msg := fmt.Sprintf("File with name \"%s\" already exists", req.Name)
		log.Warn(msg)
		return nil, nil, fmt.Errorf(msg)
	}

	newDescr, err := repo.Create(d.descriptor.GetInode(), req.Name, db.DT_File, db.DescriptorAttrs{
		GID:        req.Gid,
		UID:        req.Uid,
		Permission: "0644",
		Size:       0,
	})
	if err != nil {
		log.Errorf("Error creating file. Error: %s", err.Error())
		return nil, nil, err
	}

	file := File{
		descriptor: newDescr,
		fs:         d.fs,
	}
	handle := FileHandle{file: &file}

	return &file, &handle, nil
}

var _ = fuseFS.NodeRemover(&Dir{})

func (d *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	log.Infof("Removing entry %s in %s[%d]", req.Name, d.descriptor.GetName(), d.descriptor.GetInode())

	repo := d.fs.RepositoryRegistry.GetDescriptorRepository()
	return repo.RemoveByName(d.descriptor.GetInode(), req.Name)
}
