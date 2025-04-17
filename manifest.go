// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package quelldb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang/snappy"
	"github.com/thirashapw/quelldb/constants"
	"github.com/thirashapw/quelldb/utils"
)

// EncodeManifest encodes SSStorage names with binary format, nappy, ptional encryption
func EncodeManifest(ssts []string, key []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, int32(len(ssts)))
	for _, name := range ssts {
		binary.Write(buf, binary.LittleEndian, int32(len(name)))
		buf.Write([]byte(name))
	}
	compressed := snappy.Encode(nil, buf.Bytes())
	if key != nil {
		return utils.Encrypt(compressed, key)
	}
	return compressed, nil
}

// DecodeManifest decodes manifest data to extract SSStorage names
func DecodeManifest(data []byte, key []byte) ([]string, error) {
	if key != nil {
		var err error
		data, err = utils.Decrypt(data, key)
		if err != nil {
			return nil, err
		}
	}
	decoded, err := snappy.Decode(nil, data)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(decoded)

	var count int32
	binary.Read(buf, binary.LittleEndian, &count)

	ssts := make([]string, 0, count)
	for i := 0; i < int(count); i++ {
		var nameLen int32
		binary.Read(buf, binary.LittleEndian, &nameLen)
		name := make([]byte, nameLen)
		buf.Read(name)
		ssts = append(ssts, string(name))
	}
	return ssts, nil
}

// SaveManifest writes a new numbered manifest and updates CURRENT
func SaveManifest(basePath string, ssts []string, key []byte) error {
	// determine next manifest ID
	nextID, err := nextManifestID(basePath)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s-%05d%s", constants.MANIFEST_FILE_PREFIX, nextID, constants.MANIFEST_FILE_SUFFIX)
	fullPath := filepath.Join(basePath, filename)

	// write new manifest
	data, err := EncodeManifest(ssts, key)
	if err != nil {
		return err
	}
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return err
	}

	// update CURRENT pointer
	currentPath := filepath.Join(basePath, constants.CURRENT_MANIFEST_FILE)
	return os.WriteFile(currentPath, []byte(filename), 0644)
}

// internal: gets next manifest ID
func nextManifestID(basePath string) (int, error) {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return 1, err
	}
	max := 0
	for _, f := range files {
		name := f.Name()
		if len(name) >= 16 && name[:8] == constants.MANIFEST_FILE_PREFIX {
			numStr := name[9:14]
			num, err := strconv.Atoi(numStr)
			if err == nil && num > max {
				max = num
			}
		}
	}
	return max + 1, nil
}
