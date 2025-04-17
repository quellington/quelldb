package base

import (
	"crypto/sha256"
	"encoding/binary"
	"os"
)

type BloomFilter struct {
	bits []byte
	size uint32
	k    uint8
}

// NewBloomFilter creates a new Bloom filter with the specified size and number of hash functions.
// The size is the number of bits in the filter, and k is the number of hash functions to use.
// The size should be a multiple of 8, as the bits are stored in bytes.
// The filter is initialized with all bits set to 0.
// The filter uses SHA-256 to generate the hash values for the keys.
func ApplyNewBloomFilter(size uint32, hashCount uint8) *BloomFilter {
	return &BloomFilter{
		bits: make([]byte, size/8+1),
		size: size,
		k:    hashCount,
	}
}

func (bf *BloomFilter) Add(key string) {
	hashes := bf.getHashes(key)
	for _, h := range hashes {
		bf.setBit(h % bf.size)
	}
}

func (bf *BloomFilter) Test(key string) bool {
	hashes := bf.getHashes(key)
	for _, h := range hashes {
		if !bf.getBit(h % bf.size) {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) getHashes(key string) []uint32 {
	h := sha256.Sum256([]byte(key))
	hashes := make([]uint32, bf.k)
	for i := uint8(0); i < bf.k; i++ {
		start := i * 4
		hashes[i] = binary.LittleEndian.Uint32(h[start : start+4])
	}
	return hashes
}

func (bf *BloomFilter) setBit(pos uint32) {
	byteIndex := pos / 8
	bit := pos % 8
	bf.bits[byteIndex] |= 1 << bit
}

func (bf *BloomFilter) getBit(pos uint32) bool {
	byteIndex := pos / 8
	bit := pos % 8
	return (bf.bits[byteIndex] & (1 << bit)) != 0
}

func (bf *BloomFilter) Bytes() []byte {
	return bf.bits
}

func (bf *BloomFilter) Load(data []byte) {
	copy(bf.bits, data)
}

func saveBloomFilter(filter *BloomFilter, path string) error {
	return os.WriteFile(path, filter.Bytes(), 0644)
}

// Parameters:
// path: The path to the file where the Bloom filter will be saved.
// size: The size of the Bloom filter in bits.
// hashCount: The number of hash functions to use.
//
// Returns:
// A pointer to the BloomFilter object.
// An error if the file could not be saved.
func LoadBloomFilter(path string, size uint32, hashCount uint8) (*BloomFilter, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	filter := ApplyNewBloomFilter(size, hashCount)
	filter.Load(data)
	return filter, nil
}
