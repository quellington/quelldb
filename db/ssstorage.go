// sorted string storage
package db

import (
	"encoding/binary"
	"os"

	"github.com/golang/snappy"
	"github.com/thirashapw/quelldb/utils"
)

const xorKey byte = 0xAB

func WriteSSStorage(path string, data map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for k, v := range data {
		kb := snappy.Encode(nil, []byte(k))
		vb := snappy.Encode(nil, []byte(v))

		kb = utils.XorMask(kb, xorKey)
		vb = utils.XorMask(vb, xorKey)

		binary.Write(file, binary.LittleEndian, int32(len(kb)))
		file.Write(kb)

		binary.Write(file, binary.LittleEndian, int32(len(vb)))
		file.Write(vb)
	}
	return nil
}

func ReadSSStorage(path string) (map[string]string, error) {
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
		file.Read(kb)
		kb = utils.XorMask(kb, xorKey)
		keyDecoded, err := snappy.Decode(nil, kb)
		if err != nil {
			return nil, err
		}

		var vLen int32
		binary.Read(file, binary.LittleEndian, &vLen)
		vb := make([]byte, vLen)
		file.Read(vb)
		vb = utils.XorMask(vb, xorKey)
		valDecoded, err := snappy.Decode(nil, vb)
		if err != nil {
			return nil, err
		}

		result[string(keyDecoded)] = string(valDecoded)
	}
	return result, nil
}
