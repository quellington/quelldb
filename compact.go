// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package quelldb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thirashapw/quelldb/base"
	"github.com/thirashapw/quelldb/constants"
	"github.com/thirashapw/quelldb/utils"
)

// Compact merges multiple SSStorage into a single one.
// It reads all SSStorages in the base path, merges them into a single map,
// and writes the merged data into a new SSStorage file.
// The old SSStorage files are deleted after the merge.
func (db *DB) Compact() error {
	files, err := os.ReadDir(db.basePath)
	if err != nil {
		return err
	}

	merged := make(map[string]string)
	var sstPaths []string

	for _, f := range files {
		if strings.HasPrefix(f.Name(), constants.SSS_PREFIX) && strings.HasSuffix(f.Name(), constants.SSS_SUFFIX) {
			path := filepath.Join(db.basePath, f.Name())
			data, err := base.ReadSSStorage(path, db.key)
			if err != nil {
				return err
			}
			for k, v := range data {
				merged[k] = v
			}
			sstPaths = append(sstPaths, path)
		}
	}

	// if there are less than 2 SSStorages, no need to merge
	if uint(len(sstPaths)) < db.compactLimit {
		return nil
	}

	// write merged SSStorage
	id, _ := utils.NextSSSID(db.basePath)
	newPath := filepath.Join(db.basePath, fmt.Sprintf(constants.SSS_PREFIX+"%05d"+constants.SSS_SUFFIX, id))
	if err := base.WriteSSStorage(newPath, merged, db.key); err != nil {
		return err
	}

	// remove old SSStorage
	for _, p := range sstPaths {
		os.Remove(p)
	}

	return nil
}
