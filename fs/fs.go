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

package fs

import (
	"bazil.org/fuse"
	fuseFS "bazil.org/fuse/fs"
	"fmt"
	"github.com/kos-v/dbunderfs/db"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
)

type RootNotFoundError struct{ err string }

func (err *RootNotFoundError) Error() string {
	return fmt.Sprintf("Root \"%s\" not found.", err.err)
}

func Mount(point string, dbFactory db.DBFactory) error {
	conn, err := fuse.Mount(point)
	if err != nil {
		return err
	}
	defer conn.Close()

	filesys := &FS{
		DBFactory: dbFactory,
	}

	if err := fuseFS.Serve(conn, filesys); err != nil {
		return err
	}

	return nil
}

type FS struct {
	DBFactory db.DBFactory
}

func (f *FS) Root() (fuseFS.Node, error) {
	log.Infof("Finding root %s", db.RootName)

	descr, err := f.DBFactory.CreateDescriptorRepository().FindRoot()
	if err != nil {
		log.Warnf("Root finding error: %s", err.Error())
		return nil, err
	}
	if descr == nil {
		log.Warnf("Root not found")
		return nil, &RootNotFoundError{err: db.RootName}
	}

	return &Dir{
		descriptor: descr,
		fs:         f,
	}, nil
}

var _ fuseFS.FS = (*FS)(nil)

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

	repo := d.fs.DBFactory.CreateDescriptorRepository()
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

	collection, err := d.fs.DBFactory.CreateDescriptorRepository().FindChildrenByInode(d.descriptor.GetInode())
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

	repo := d.fs.DBFactory.CreateDescriptorRepository()
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

	repo := d.fs.DBFactory.CreateDescriptorRepository()
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

	return &file, &file, nil
}

var _ = fuseFS.NodeRemover(&Dir{})

func (d *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	log.Infof("Removing entry %s in %s[%d]", req.Name, d.descriptor.GetName(), d.descriptor.GetInode())

	repo := d.fs.DBFactory.CreateDescriptorRepository()
	return repo.RemoveByName(d.descriptor.GetInode(), req.Name)
}

type File struct {
	descriptor db.DescriptorInterface
	fs         *FS
}

var _ fuseFS.Node = (*File)(nil)

func (f *File) Attr(ctx context.Context, attr *fuse.Attr) error {
	log.Infof("Reads file attrs of %d:%s", f.descriptor.GetInode(), db.RootName)

	attr.Inode = uint64(f.descriptor.GetInode())
	attr.Mode = f.descriptor.GetPermission()
	attr.Uid = f.descriptor.GetUID()
	attr.Gid = f.descriptor.GetGID()
	attr.Size = f.descriptor.GetSize()

	return nil
}

var _ = fuseFS.NodeOpener(&File{})

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fuseFS.Handle, error) {
	log.Infof("Opening file %d:%s", f.descriptor.GetInode(), f.descriptor.GetName())

	//resp.Flags |= fuse.OpenNonSeekable
	resp.Flags |= fuse.OpenDirectIO

	return &FileHandle{file: f}, nil
}

type FileHandle struct {
	file *File
}

var _ fuseFS.Handle = (*FileHandle)(nil)

var _ fuseFS.HandleReleaser = (*FileHandle)(nil)

func (fh *FileHandle) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	return nil
}

var _ = fuseFS.HandleReader(&FileHandle{})

func (fh *FileHandle) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	descr := fh.file.descriptor
	log.Infof("Reading file %d:%s", descr.GetInode(), descr.GetName())

	repo := fh.file.fs.DBFactory.CreateDataBlockRepository()
	dataBlock, err := repo.FindFirst(descr)
	if err != nil {
		log.Errorf("Error: %s", err.Error())
		return err
	}
	if dataBlock == nil {
		msg := fmt.Sprintf("Data block for descriptor %d:%s not found", descr.GetInode(), descr.GetName())
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	data := *dataBlock.GetData()
	if req.Offset >= int64(len(data)) {
		data = nil
	} else {
		data = data[req.Offset:]
	}
	if len(data) > req.Size {
		data = data[:req.Size]
	}
	n := copy(resp.Data[:req.Size], data)
	resp.Data = resp.Data[:n]

	return nil
}

var _ = fuseFS.HandleWriter(&FileHandle{})

func (fh *FileHandle) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	descr := fh.file.descriptor
	log.Infof("Write file %d:%s", descr.GetInode(), descr.GetName())

	repo := fh.file.fs.DBFactory.CreateDataBlockRepository()
	dataBlock, err := repo.FindFirst(descr)
	if err != nil {
		log.Errorf("Error: %s", err.Error())
		return err
	}

	addedSize := dataBlock.Add(uint64(req.Offset), &req.Data)

	err = repo.Write(descr, dataBlock.GetData())
	if err != nil {
		log.Errorf("Error write to %s[%d] file. Error: %s", descr.GetName(), descr.GetInode(), err)
		return err
	}

	resp.Size = addedSize

	return nil
}
