// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

// sorted string storage
package base

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/golang/snappy"
	"github.com/quellington/quelldb/constants"
	"github.com/quellington/quelldb/utils"
)

// WriteSSStorage writes a map of strings to a file in a sorted string storage format.
// Each key-value pair is compressed using snappy and optionally encrypted.
// The keys and values are prefixed with their lengths to allow for easy reading.
// The file is created if it doesn't exist, and overwritten if it does.
// The path parameter specifies the file location, and the key parameter is used for encryption.
// If the key is nil, the data will be stored unencrypted.
func WriteSSStorage(path string, data map[string]string, key []byte) (string, string, error) {

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	minKey := keys[0]
	maxKey := keys[len(keys)-1]

	file, err := os.Create(path)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	offsets := make(map[string]int64)

	filter := ApplyNewBloomFilter(constants.BOOM_BIT_SIZE, constants.BOOM_HASH_COUNT)

	for k, v := range data {

		// get current byte offset
		pos, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			return "", "", err
		}
		offsets[k] = pos

		filter.Add(k)
		kb := snappy.Encode(nil, []byte(k))
		vb := snappy.Encode(nil, []byte(v))

		if key != nil {
			kb, err = utils.Encrypt(kb, key)
			if err != nil {
				return "", "", err
			}
			vb, err = utils.Encrypt(vb, key)
			if err != nil {
				return "", "", err
			}
		}

		binary.Write(file, binary.LittleEndian, int32(len(kb)))
		file.Write(kb)
		binary.Write(file, binary.LittleEndian, int32(len(vb)))
		file.Write(vb)
	}

	// serialize the index map
	indexBytes, err := json.Marshal(offsets)
	if err != nil {
		return "", "", err
	}
	_, err = file.Write(indexBytes)
	if err != nil {
		return "", "", err
	}

	// write the length of index
	binary.Write(file, binary.LittleEndian, int32(len(indexBytes)))

	file.Write([]byte(constants.INDEX_FOOTER_NAME))

	// Save bloom filter
	err = saveBloomFilter(filter, path+constants.SSS_BOOM_FILTER_SUFFIX)
	if err != nil {
		return "", "", err
	}

	return minKey, maxKey, nil
}

// ReadSSStorage reads a sorted string storage file and returns a map of strings.
// Each key-value pair is read from the file, and the values are decompressed using snappy.
// If the key parameter is provided, the data will be decrypted using the key.
// If the key is nil, the data will be read unencrypted.
// Function also handles the case where the file is too small to contain a valid index.
// Direct file seek instead of using a buffered reader to avoid memory overhead.
func ReadSSStorage(path string, key []byte) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// load index footer
	stat, _ := file.Stat()
	if stat.Size() < 8 {
		return nil, fmt.Errorf("file too small to contain index")
	}

	// Seek to index footer
	_, err = file.Seek(-8, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	var indexLen int32
	err = binary.Read(file, binary.LittleEndian, &indexLen)
	if err != nil {
		return nil, err
	}

	footer := make([]byte, 4)
	_, err = file.Read(footer)
	if err != nil || string(footer) != constants.INDEX_FOOTER_NAME {
		return nil, fmt.Errorf("invalid SSStorage format: missing footer")
	}

	// seek to index and read it
	_, err = file.Seek(-(int64(indexLen) + 8), io.SeekEnd)
	if err != nil {
		return nil, err
	}

	indexBytes := make([]byte, indexLen)
	_, err = file.Read(indexBytes)
	if err != nil {
		return nil, err
	}

	offsetMap := map[string]int64{}
	if err := json.Unmarshal(indexBytes, &offsetMap); err != nil {
		return nil, err
	}

	// load all keys based on offsets
	result := make(map[string]string)
	for _, offset := range offsetMap {
		_, err := file.Seek(offset, io.SeekStart)
		if err != nil {
			return nil, err
		}

		var kLen int32
		err = binary.Read(file, binary.LittleEndian, &kLen)
		if err != nil {
			return nil, err
		}

		kb := make([]byte, kLen)
		_, err = file.Read(kb)
		if err != nil {
			return nil, err
		}
		if key != nil {
			kb, err = utils.Decrypt(kb, key)
			if err != nil {
				return nil, err
			}
		}
		decodedKey, err := snappy.Decode(nil, kb)
		if err != nil {
			return nil, err
		}

		var vLen int32
		err = binary.Read(file, binary.LittleEndian, &vLen)
		if err != nil {
			return nil, err
		}
		vb := make([]byte, vLen)
		_, err = file.Read(vb)
		if err != nil {
			return nil, err
		}
		if key != nil {
			vb, err = utils.Decrypt(vb, key)
			if err != nil {
				return nil, err
			}
		}
		valDecoded, err := snappy.Decode(nil, vb)
		if err != nil {
			return nil, err
		}

		result[string(decodedKey)] = string(valDecoded)
	}

	return result, nil
}
