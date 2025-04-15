// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

// write ahead log (WAL) implementation
package base

import (
	"fmt"
	"os"
)

type WAL struct {
	file *os.File
}

// NewWAL creates a new Write Ahead Log (WAL) at the specified path.
// It opens the file for appending and creates it if it doesn't exist.
func NewWAL(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{file: f}, nil
}


// Write appends a new log entry to the WAL.
// Each entry consists of an operation (op), a key, and a value.
// The entry is formatted as "op|key|value\n".
func (w *WAL) Write(op, key, value string) error {
	line := fmt.Sprintf("%s|%s|%s\n", op, key, value)
	_, err := w.file.WriteString(line)
	return err
}


// Read reads all log entries from the WAL.
func (w *WAL) Close() error {
	return w.file.Close()
}
