// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package constants

const (
	LOG_FILE                  = "00000.log"
	SSS_MERGE_FILE_NAME       = "sss-merged"
	SSS_PREFIX                = "sss-"
	SSS_SUFFIX                = ".qldb"
	SSS_BOOM_FILTER_SUFFIX    = ".filter"
	INDEX_FOOTER_NAME         = "QIDX"
	SSS_COMPACT_DEFAULT_LIMIT = 10
	BOOM_BIT_SIZE             = 8000
	BOOM_HASH_COUNT           = 4

	// KEY
	PUT    = "PUT"
	DELETE = "DEL"
	GET    = "GET"
	ALL    = "ALL"


	// MANIFEST
	CURRENT_MANIFEST_FILE = "CURRENT"
	MANIFEST_FILE_PREFIX  = "MANIFEST"
	MANIFEST_FILE_SUFFIX  = ".qmf"
)
