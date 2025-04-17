// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package quelldb

import (
	"bytes"
	"encoding/binary"

	"github.com/golang/snappy"
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
