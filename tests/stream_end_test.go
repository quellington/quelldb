// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.
package main

import (
	"testing"
	"time"

	"github.com/thirashapw/quelldb"
)

func TestStreamEnd(t *testing.T) {
	db, _ := quelldb.Open("data", nil)

	id := db.Subscribe(func(e quelldb.ChangeEvent) {
		t.Logf("[Event] %s - %s = %s", e.Type, e.Key, e.Value)
	})

	db.Put("user:1", "alice")
	db.Delete("user:1")

	db.Unsubscribe(id)

	db.Put("user:2", "bob")

	time.Sleep(1000 * time.Millisecond)
}
