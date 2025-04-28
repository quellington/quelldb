# QuellDB

QuellDB is a lightweight, embeddable, high-performance key-value store written in Go.  
Built on a Log-Structured Merge Tree (LSM Tree) architecture, QuellDB provides:

- Fast in-memory writes
- Durable disk-based persistence
- Optional AES-256 encryption
- Binary-safe, compressed SSStorages using Snappy
- Range-aware compaction with manifest tracking
- Simple, pluggable Go API

---

## Features

- In-memory `MemStorage` with thread-safe access
- Persistent `SSStorages` (Sorted String Storages)
- Optional AES-256 encryption with GCM mode
- Snappy compression by default (even without encryption)
- Write-Ahead Log (WAL) for durability before flush
- Bloom filter support for efficient lookups
- TTL (Time-To-Live) support for expiring keys
- Batch writes via `PutBatch()`
- Key iteration via `Iterator()`, with prefix filters
- Versioned manifest system
- Range-aware SSStorage compaction based on overlapping key ranges

---

## Installation

```bash
go get github.com/quellington/quelldb
```

---

## Quick Start
Snappy compression (Default)

```bash
package main

import (
    "fmt"
    "github.com/quellington/quelldb"
)

func main() {
    store, err := quelldb.Open("data", nil)
    if err != nil {
        panic(err)
    }

    store.Put("username", "john")
    store.Flush()

    val, _ := store.Get("username")
    fmt.Println("Value:", val)
}

```


With AES-256 Encryption (If needed)

```bash
package main

import (
    "fmt"
    "github.com/quellington/quelldb"
)

func main() {
    key := []byte("thisis32byteslongthisis32byteslo") // 32 bytes

    store, err := quelldb.Open("securedata", &quelldb.Options{
        EncryptionKey: key,
    })
    if err != nil {
        panic(err)
    }

    store.Put("user:123", `{"id":"123","username":"john","email":"john@example.io","age":50}`)
    store.Flush()

    val, _ := store.Get("user:123")
    fmt.Println("Decrypted:", val)
}

```

TTL Support

```bash
store.PutTTL("temp:session", "expires-soon", 10*time.Second)
```


Batch Writes

```bash
store.PutBatch(map[string]string{
    "user:1": "alice",
    "user:2": "bob",
})
```

Prefix Iteration

```bash
it := store.PrefixIterator("user:")
for it.Next() {
    fmt.Println(it.Key(), it.Value())
}

```


### Encryption Logic
If you provide a 32-byte EncryptionKey, QuellDB will:

- Compress values with Snappy
- Encrypt them using AES-256 GCM mode
- Without a key, values are still compressed and unreadable to humans, but not encrypted


## API Reference

| Function       | Description                                      |
|----------------|--------------------------------------------------|
| `Open()`    | Initializes a database at given path             |
| `Put(key, val)`| Writes data into memory and WAL                  |
| `Get(key)`     | Retrieves value from memory or SSStorages        |
| `Delete(key)`     | Deletes a key from memory and appends `DEL` to WAL        |
| `Flush()`      | Persists current MemStorage to a new SSStorage   |
| `PutBatch(map[string]string)`      | PeWrites multiple key-value pairs in one WAL flush   |
| `PutTTL(key, val, ttl)`      | Writes a key with an expiration duration   |
| `Iterator()`      | Iterates all sorted keys from memory   |
| `PrefixIterator(p)`      | Iterates sorted keys with the given prefix   |
| `Compact(p)`      | Compacts overlapping SSStorage into a single one   |
| `Subscribe(func(ChangeEvent)) int`      | Registers a live event handler and returns a handler ID   |

MIT License © 2025 The QuellDB Authors