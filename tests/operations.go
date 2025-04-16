// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

// TEST
package main

import (
	"encoding/json"
	"fmt"

	"github.com/thirashapw/quelldb"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

func SaveUser(db *quelldb.DB, user User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return db.Put("user:"+user.ID, string(userBytes))
}
func LoadUser(db *quelldb.DB, id string) (*User, error) {
	data, ok := db.Get("user:" + id)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	var user User
	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func main() {
	store, err := quelldb.Open("data", nil)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	store.Put("foo", "bar")
	store.Put("hello", "world")
	store.Put("hedsadllo", "world")
	store.Put("heldsadsalo", "world")
	store.Put("hedsadllo", "world")
	store.Put("heldsdsalo", "world")
	store.Put("heldsadlo", "wodsdrld")
	store.Put("heldsadlo", "worldsdd")
	store.Put("heldsadlo", "world")

	val, _ := store.Get("heldsadlo")
	fmt.Println("Value of foo:", val)

	u := User{
		ID:       "123",
		Username: "thirasha",
		Email:    "t@crypto.io",
		Age:      50,
	}
	SaveUser(store, u)
	store.Flush()
	loadedUser, _ := LoadUser(store, "123")
	fmt.Println("Username:", loadedUser.Username)

	// AES

	// key := []byte("thisis32byteslongthisis32byteslo")

	// store, err := quelldb.Open("securedata", &quelldb.Options{
	// 	EncryptionKey: key,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// store.Put("email", "user@example.com")
	// store.Flush()

	// val, ok := store.Get("email")
	// if ok {
	// 	fmt.Println("Encrypted Value:", val)
	// }

	// store.Flush()
}
