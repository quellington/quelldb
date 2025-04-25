// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package quelldb

type ChangeEvent struct {
	Type  string
	Key   string
	Value string
}

// Subscribe allows you to register a callback function that will be called
// whenever a change event occurs in the database.
// The callback function receives a ChangeEvent struct containing the type of change,
// the key, and the value.
func (db *DB) Subscribe(handler func(ChangeEvent)) int {
	db.subLock.Lock()
	defer db.subLock.Unlock()
	if db.subscribers == nil {
		db.subscribers = make(map[int]func(ChangeEvent))
	}

	id := db.nextSubID
	db.nextSubID++
	db.subscribers[id] = handler

	return id
}

// Unsubscribe removes a subscriber from the database.
// The subscriber is identified by its ID, which is returned when the subscriber was added.
func (db *DB) Unsubscribe(id int) {
	db.subLock.Lock()
	defer db.subLock.Unlock()

	delete(db.subscribers, id)
}

// Subscribe adds a new subscriber to the database.
// The subscriber is a function that takes a ChangeEvent as an argument.
// The subscriber will be notified of changes to the database, such as PUT or DELETE operations.
func (db *DB) publish(event ChangeEvent) {
	db.subLock.RLock()
	defer db.subLock.RUnlock()
	for _, handler := range db.subscribers {
		go handler(event)
	}
}
