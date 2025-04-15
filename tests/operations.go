// TEST
package main

import (
	"encoding/json"
	"fmt"

	"github.com/thirashapw/quelldb/db"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

func SaveUser(db *db.DB, user User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return db.Put("user:"+user.ID, string(userBytes))
}
func LoadUser(db *db.DB, id string) (*User, error) {
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
	store, err := db.Open("data", nil)
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

	loadedUser, _ := LoadUser(store, "123")
	fmt.Println("Username:", loadedUser.Username)

	store.Flush()
}
