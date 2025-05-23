// TEST
package tests

import (
	"encoding/json"
	"testing"

	"github.com/quellington/quelldb"
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
	data, gErr := db.Get("user:" + id)
	if gErr != nil {
		return nil, gErr
	}
	var user User
	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func TestMain(t *testing.T) {
	store, err := quelldb.Open("data", &quelldb.Options{
		CompactLimit:  10,
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
	t.Logf("Value of foo: %s", val)

	// u := User{
	// 	ID:       "123",
	// 	Username: "thirasha",
	// 	Email:    "t@crypto.io",
	// 	Age:      50,
	// }
	// SaveUser(store, u)
	// store.Flush()
	// store.Compact()
	// loadedUser, _ := LoadUser(store, "123")
	// fmt.Println("Username:", loadedUser)

	// // ----

	// users := map[string]string{}

	// u1 := User{Username: "john"}
	// u2 := User{Username: "sarah"}
	// u3 := User{Username: "mike"}

	// b1, _ := json.Marshal(u1)
	// b2, _ := json.Marshal(u2)
	// b3, _ := json.Marshal(u3)

	// users["user:101"] = string(b1)
	// users["user:102"] = string(b2)
	// users["user:103"] = string(b3)

	// err = store.PutBatch(users)
	// if err != nil {
	// 	panic(err)
	// }

	// store.Flush()
	// store.Compact()

	// data, gErr := store.Get("user:102")
	// if gErr != nil || data == "" {
	// 	log.Fatal(gErr.Error())
	// }

	// var user User
	// if err := json.Unmarshal([]byte(data), &user); err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Username2:", user.Username)

	// it := store.Iterator()

	// for it.Next() {
	// 	fmt.Println(it.Key(), it.Value())
	// }

	it := store.PrefixIterator("user:")
	for it.Next() {
		t.Logf(it.Key(), it.Value())
	}

	// ----

	// -- ttl --

	// store.PutTTL("test:token", "abc123", 1*time.Second)

	// time.Sleep(4 * time.Second)

	// val, ok = store.Get("test:token")
	// fmt.Println("Value of test:token:", val, ok)

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
