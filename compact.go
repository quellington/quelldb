// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package quelldb

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thirashapw/quelldb/base"
	"github.com/thirashapw/quelldb/constants"
	"github.com/thirashapw/quelldb/utils"
)

// Compact merges multiple SSStorage into a single one.
// It reads all SSStorages in the base path, merges them into a single map,
// and writes the merged data into a new SSStorage file.
// The old SSStorage files are deleted after the merge.
func (db *DB) Compact() error {

	if len(db.manifestSSSs) < int(db.compactLimit) {
		return nil
	}

	var toCompact []SSSMeta
	basep := db.manifestSSSs[0]
	toCompact = append(toCompact, basep)
	// select overlapping SSTs for compaction
	for _, sst := range db.manifestSSSs[1:] {
		if overlapsAny(sst, toCompact) {
			toCompact = append(toCompact, sst)
		}
	}

	// if not enough overlapping SSTs to trigger compaction
	if uint(len(toCompact)) < db.compactLimit {
		return nil
	}

	merged := make(map[string]string)

	for _, f := range toCompact {
		fullPath := filepath.Join(db.basePath, f.Filename)
		data, err := base.ReadSSStorage(fullPath, db.key)
		if err != nil {
			return err
		}
		for k, v := range data {
			merged[k] = v
		}
		os.Remove(fullPath)
		os.Remove(fullPath + constants.SSS_BOOM_FILTER_SUFFIX)
	}

	// if there are less than 2 SSStorages, no need to merge
	// if uint(len(sstPaths)) < db.compactLimit {
	// 	return nil
	// }

	// write merged SSStorage
	id, _ := utils.NextSSSID(db.basePath)
	newSSSFile := fmt.Sprintf(constants.SSS_PREFIX+"%05d"+constants.SSS_SUFFIX, id)
	newPath := filepath.Join(db.basePath, newSSSFile)
	minKey, maxKey, err := base.WriteSSStorage(newPath, merged, db.key)
	if err != nil {
		return err
	}

	// update manifest
	db.manifestSSSs = removeCompactedSSSs(db.manifestSSSs, toCompact)
	db.manifestSSSs = append(db.manifestSSSs, SSSMeta{
		Filename: newSSSFile,
		MinKey:   minKey,
		MaxKey:   maxKey,
	})

	return SaveManifest(db.basePath, db.manifestSSSs, db.key)
}
