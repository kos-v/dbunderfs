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
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

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

	repo := fh.file.fs.RepositoryRegistry.GetDataBlockRepository()
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

	repo := fh.file.fs.RepositoryRegistry.GetDataBlockRepository()
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
