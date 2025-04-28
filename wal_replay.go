// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package quelldb

import (
	"bufio"
	"os"
	"strings"

	"github.com/quellington/quelldb/constants"
)

// Put stores a key-value pair in the database.
// It first stores the pair in memory and then writes it to the WAL.
// The function takes a key and a value as parameters.
func (db *DB) replayWAL(path string) error {
	file, err := os.Open(path)
	if err != nil {

		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 3)

		// skip empty lines or lines that don't have exact parts
		if len(parts) != 3 {
			continue
		}
		op, key, val := parts[0], parts[1], parts[2]
		switch op {
		case constants.GET:
			db.memStorage.Put(key, val)
		case constants.DELETE:
			db.memStorage.Delete(key)
		}
	}
	return scanner.Err()
}
