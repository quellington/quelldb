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
	"strings"

	"github.com/golang/snappy"
	"github.com/quellington/quelldb/constants"
	"github.com/quellington/quelldb/utils"
)

type SSSMeta struct {
	Filename string
	MinKey   string
	MaxKey   string
}

// write string as [len][bytes]
func writeString(buf *bytes.Buffer, s string) {
	binary.Write(buf, binary.LittleEndian, int32(len(s)))
	buf.Write([]byte(s))
}

// read string as [len][bytes]
func readString(buf *bytes.Reader) string {
	var n int32
	binary.Read(buf, binary.LittleEndian, &n)
	b := make([]byte, n)
	buf.Read(b)
	return string(b)
}

// EncodeManifest encodes SSStorage names with binary format, nappy, ptional encryption
func EncodeManifest(ssts []SSSMeta, key []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, int32(len(ssts)))
	for _, s := range ssts {
		writeString(buf, s.Filename)
		writeString(buf, s.MinKey)
		writeString(buf, s.MaxKey)
	}
	compressed := snappy.Encode(nil, buf.Bytes())
	if key != nil {
		return utils.Encrypt(compressed, key)
	}
	return compressed, nil
}

// DecodeManifest decodes manifest data to extract SSStorage names
func DecodeManifest(data []byte, key []byte) ([]SSSMeta, error) {
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

	ssts := make([]SSSMeta, 0, count)
	for i := 0; i < int(count); i++ {
		meta := SSSMeta{
			Filename: readString(buf),
			MinKey:   readString(buf),
			MaxKey:   readString(buf),
		}
		ssts = append(ssts, meta)
	}
	return ssts, nil
}

// SaveManifest writes a new numbered manifest and updates CURRENT
func SaveManifest(basePath string, ssts []SSSMeta, key []byte) error {
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
	if err := os.WriteFile(currentPath, []byte(filename), 0644); err != nil {
		return err
	}

	// manifest cleanup
	files, _ := os.ReadDir(basePath)
	for _, f := range files {
		name := f.Name()
		if strings.HasPrefix(name, constants.MANIFEST_FILE_PREFIX+"-") && name != filename {
			os.Remove(filepath.Join(basePath, name))
		}
	}

	return nil
}

// LoadManifest reads CURRENT, then loads the correct numbered manifest
func LoadManifest(basePath string, key []byte) ([]SSSMeta, error) {
	currentPath := filepath.Join(basePath, constants.CURRENT_MANIFEST_FILE)
	data, err := os.ReadFile(currentPath)
	if err != nil {
		if os.IsNotExist(err) {

			// CURRENT file not found, return empty manifest (fresh storage)
			return []SSSMeta{}, nil
		}
		return nil, err
	}
	manifestName := string(bytes.TrimSpace(data))
	manifestPath := filepath.Join(basePath, manifestName)
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	return DecodeManifest(manifestData, key)
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
		if len(name) >= 16 && strings.HasPrefix(name, constants.MANIFEST_FILE_PREFIX) {
			numStr := name[len(constants.MANIFEST_FILE_PREFIX)+1 : len(constants.MANIFEST_FILE_PREFIX)+6]
			num, err := strconv.Atoi(numStr)
			if err == nil && num > max {
				max = num
			}
		}
	}
	return max + 1, nil
}

func overlapsAny(a SSSMeta, group []SSSMeta) bool {
	for _, b := range group {
		if !(a.MaxKey < b.MinKey || a.MinKey > b.MaxKey) {
			return true
		}
	}
	return false
}

func removeCompactedSSSs(all []SSSMeta, toRemove []SSSMeta) []SSSMeta {
	removeMap := make(map[string]bool)
	for _, s := range toRemove {
		removeMap[s.Filename] = true
	}
	var result []SSSMeta
	for _, s := range all {
		if !removeMap[s.Filename] {
			result = append(result, s)
		}
	}
	return result
}
