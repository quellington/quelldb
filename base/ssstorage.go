// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

// sorted string storage
package base

import (
	"encoding/binary"
	"os"

	"github.com/golang/snappy"
	"github.com/thirashapw/quelldb/utils"
)

func WriteSSStorage(path string, data map[string]string, key []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for k, v := range data {
		kb := snappy.Encode(nil, []byte(k))
		vb := snappy.Encode(nil, []byte(v))

		if key != nil {
			kb, err = utils.Encrypt(kb, key)
			if err != nil {
				return err
			}
			vb, err = utils.Encrypt(vb, key)
			if err != nil {
				return err
			}
		}

		binary.Write(file, binary.LittleEndian, int32(len(kb)))
		file.Write(kb)
		binary.Write(file, binary.LittleEndian, int32(len(vb)))
		file.Write(vb)
	}
	return nil
}

func ReadSSStorage(path string, key []byte) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]string)
	for {
		var kLen int32
		err := binary.Read(file, binary.LittleEndian, &kLen)
		if err != nil {
			break
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
		keyDecoded, err := snappy.Decode(nil, kb)
		if err != nil {
			return nil, err
		}

		var vLen int32
		binary.Read(file, binary.LittleEndian, &vLen)
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
		result[string(keyDecoded)] = string(valDecoded)
	}
	return result, nil
}
