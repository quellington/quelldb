package benchmark

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/thirashapw/quelldb"
)

func setupDB(tb interface{ Fatalf(string, ...any) }) *quelldb.DB {
	os.RemoveAll("tmp-bench")
	db, err := quelldb.Open("tmp-bench", nil)
	if err != nil {
		tb.Fatalf("Failed to open DB: %v", err)
	}
	return db
}

func BenchmarkPut(b *testing.B) {
	db := setupDB(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		val := "value" + strconv.Itoa(i)
		db.Put(key, val)
	}
}

func BenchmarkPutWithFlush(b *testing.B) {
	db := setupDB(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "flush-key" + strconv.Itoa(i)
		val := "flush-val" + strconv.Itoa(i)
		db.Put(key, val)

		if i%1000 == 0 {
			db.Flush()
		}
	}
}

func BenchmarkGet(b *testing.B) {
	db := setupDB(b)
	// Pre-populate
	for i := 0; i < 10000; i++ {
		db.Put("getkey"+strconv.Itoa(i), "val")
	}
	db.Flush()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Get("getkey" + strconv.Itoa(i%10000))
	}
}

func BenchmarkFlush(b *testing.B) {
	db := setupDB(b)

	for i := 0; i < 100000; i++ {
		db.Put("flush"+strconv.Itoa(i), "val")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Flush()
	}
}

func BenchmarkPutWithTTL(b *testing.B) {
	db := setupDB(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.PutTTL("ttlkey"+strconv.Itoa(i), "val", 2*time.Second)
	}
}
