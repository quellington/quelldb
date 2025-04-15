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

func NewWAL(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{file: f}, nil
}

func (w *WAL) Write(op, key, value string) error {
	line := fmt.Sprintf("%s|%s|%s\n", op, key, value)
	_, err := w.file.WriteString(line)
	return err
}

func (w *WAL) Close() error {
	return w.file.Close()
}
