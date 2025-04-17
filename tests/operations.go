// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

// TEST
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	store, err := quelldb.Open("data", &quelldb.Options{
		CompactLimit:  5,
		BoomHashCount: 4,
	})
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

	val, _ := store.Get("hedsadllo")
	fmt.Println("Value of foo:", val)

	u := User{
		ID:       "123",
		Username: "thirasha",
		Email:    "t@crypto.io",
		Age:      50,
	}
	SaveUser(store, u)
	store.Flush()
	store.Compact()
	// loadedUser, _ := LoadUser(store, "1234")
	// fmt.Println("Username:", loadedUser)

	// ----

	users := map[string]string{}

	u1 := User{Username: "john"}
	u2 := User{Username: "sarah"}
	u3 := User{Username: "mike"}

	b1, _ := json.Marshal(u1)
	b2, _ := json.Marshal(u2)
	b3, _ := json.Marshal(u3)

	users["user:101"] = string(b1)
	users["user:102"] = string(b2)
	users["user:103"] = string(b3)

	err = store.PutBatch(users)
	if err != nil {
		panic(err)
	}

	store.Flush()
	store.Compact()

	data, ok := store.Get("user:102")
	if !ok || data == "" {
		log.Fatal("Key not found or empty")
	}

	var user User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		panic(err)
	}

	// fmt.Println("Username2:", user.Username)

	// it := store.Iterator()

	// for it.Next() {
	// 	fmt.Println(it.Key(), it.Value())
	// }

	it := store.PrefixIterator("user:")
	for it.Next() {
		fmt.Println(it.Key(), it.Value())
	}

	// ----

	// -- ttl --

	store.PutTTL("test:token", "abc123", 5*time.Second)

	time.Sleep(4 * time.Second)

	val, ok = store.Get("test:token")
	fmt.Println("Value of test:token:", val, ok)

	// -- ttl --

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
