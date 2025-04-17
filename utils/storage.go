// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package utils

import (
	"os"
	"strconv"
	"strings"

	"github.com/thirashapw/quelldb/constants"
)

// nextSSSID generates the next SSStorage ID based on existing files in the base path.
func NextSSSID(basePath string) (int, error) {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return 0, err
	}
	maxID := 0
	for _, f := range files {
		if strings.HasPrefix(f.Name(), constants.SSS_PREFIX) {
			idStr := strings.TrimSuffix(strings.TrimPrefix(f.Name(), constants.SSS_PREFIX), constants.SSS_SUFFIX)
			id, _ := strconv.Atoi(idStr)
			if id > maxID {
				maxID = id
			}
		}
	}
	return maxID + 1, nil
}
