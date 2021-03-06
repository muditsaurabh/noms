// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package nbs

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/attic-labs/noms/go/d"
)

type fsTablePersister struct {
	dir string
}

func (ftp fsTablePersister) Compact(mt *memTable, haver chunkReader) chunkSource {
	tempName, name, chunkCount := func() (string, addr, uint32) {
		var tempName string
		name, data, chunkCount := mt.write(haver)
		if chunkCount > 0 {
			temp, err := ioutil.TempFile(ftp.dir, "nbs_table_")
			d.PanicIfError(err)
			defer checkClose(temp)
			io.Copy(temp, bytes.NewReader(data))
			tempName = temp.Name()
		}
		return tempName, name, chunkCount
	}()
	if chunkCount == 0 {
		return emptyChunkSource{}
	}
	err := os.Rename(tempName, filepath.Join(ftp.dir, name.String()))
	d.PanicIfError(err)
	return newMmapTableReader(ftp.dir, name, chunkCount)
}

func (ftp fsTablePersister) Open(name addr, chunkCount uint32) chunkSource {
	return newMmapTableReader(ftp.dir, name, chunkCount)
}
