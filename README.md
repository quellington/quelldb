# QuellDB

QuellDB is a lightweight, embeddable, high-performance key-value store written in Go.  
Built on a Log-Structured Merge Tree (LSM Tree) architecture, QuellDB provides:

- Fast in-memory writes
- Durable disk-based persistence
- Optional AES-256 encryption
- Binary-safe, compressed SSStorages using Snappy
- Simple, pluggable Go API

---

## Features

- In-memory `MemStorage` with thread-safe access
- Persistent `SSStorages` (Sorted String Storages)
- Optional AES-256 encryption with GCM mode
- Snappy compression by default (even without encryption)
- Write-Ahead Log (WAL) for durability before flush
- Zero external database dependency
- Modular architecture for extension

---

## Limitations

- Indexing and secondary indexes are not supported. All data access is key-based.
- The library is designed for embedded use only. There is no built-in client-server model or remote access.

## Installation

```bash
go get github.com/thirashapw/quelldb
```

---

## Quick Start
Snappy compression (Default)

```bash
package main

import (
    "fmt"
    "github.com/thirashapw/quelldb"
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
    "github.com/thirashapw/quelldb"
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

### Encryption Logic
If you provide a 32-byte EncryptionKey, QuellDB will:

- Compress values with Snappy
- Encrypt them using AES-256 GCM mode
- Without a key, values are still compressed and unreadable to humans, but not encrypted


## API Reference

| Function       | Description                                      |
|----------------|--------------------------------------------------|
| `quelldb.Open()`    | Initializes a database at given path             |
| `Put(key, val)`| Writes data into memory and WAL                  |
| `Get(key)`     | Retrieves value from memory or SSStorages        |
| `Flush()`      | Persists current MemStorage to a new SSStorage   |

MIT License © 2025 The QuellDB Authors


## Contributors
- [@thirashapw](https://github.com/thirashapw) – Creator and Maintainer
