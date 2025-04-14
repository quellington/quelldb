package main

import (
	"fmt"

	"github.com/thirashapw/quelldb/db"
)

func main() {
	store, err := db.Open("data")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	store.Put("foo", "bar")
	store.Put("hello", "world")

	val, _ := store.Get("foo")
	fmt.Println("Value of foo:", val)

	store.Flush()
}
