// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package tests

import (
	"testing"
	"time"

	"github.com/quellington/quelldb"
)

func TestStream(t *testing.T) {
	t.Logf("test stream starting...")
	db, _ := quelldb.Open("data", nil)
	db.Subscribe(func(e quelldb.ChangeEvent) {
		t.Logf("[Event] %s - Key: %s, Value: %s\n", e.Type, e.Key, e.Value)
	})
	t.Logf("starting operations...")
	db.Put("user:1", "alice")
	db.Put("user:2", "bob")
	db.Delete("user:1")
	t.Logf("end operations...")
	time.Sleep(500 * time.Millisecond)
}
