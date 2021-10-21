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
	"io/fs"
)

type RootNotFoundError struct{ err string }

func (err *RootNotFoundError) Error() string {
	return fmt.Sprintf("Root \"%s\" not found.", err.err)
}

func Mount(point string, repositoryRegistry db.RepositoryRegistry) error {
	conn, err := fuse.Mount(point)
	if err != nil {
		return err
	}
	defer conn.Close()

	filesys := &FS{
		RepositoryRegistry: repositoryRegistry,
	}

	if err := fuseFS.Serve(conn, filesys); err != nil {
		return err
	}

	return nil
}

type FS struct {
	RepositoryRegistry db.RepositoryRegistry
}

func (f *FS) Root() (fuseFS.Node, error) {
	log.Infof("Finding root %s", db.RootName)

	descr, err := f.RepositoryRegistry.GetDescriptorRepository().FindRoot()
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

type Permission fs.FileMode

func (p Permission) ToOctalString() string {
	return fmt.Sprintf("%#o", p)
}
