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
	"github.com/kos-v/dbunderfs/internal/db"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

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
