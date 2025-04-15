// sorted string storage
package db

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
		keyBytes := snappy.Encode(nil, []byte(k))
		valBytes := snappy.Encode(nil, []byte(v))

		if key != nil {
			keyBytes, err = utils.Encrypt(keyBytes, key)
			if err != nil {
				return err
			}
			valBytes, err = utils.Encrypt(valBytes, key)
			if err != nil {
				return err
			}
		}

		binary.Write(file, binary.LittleEndian, int32(len(keyBytes)))
		file.Write(keyBytes)
		binary.Write(file, binary.LittleEndian, int32(len(valBytes)))
		file.Write(valBytes)
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
		var keyLen int32
		err := binary.Read(file, binary.LittleEndian, &keyLen)
		if err != nil {
			break
		}

		keyBytes := make([]byte, keyLen)
		_, err = file.Read(keyBytes)
		if err != nil {
			return nil, err
		}
		if key != nil {
			keyBytes, err = utils.Decrypt(keyBytes, key)
			if err != nil {
				return nil, err
			}
		}
		keyDecoded, err := snappy.Decode(nil, keyBytes)
		if err != nil {
			return nil, err
		}

		var valLen int32
		binary.Read(file, binary.LittleEndian, &valLen)
		valBytes := make([]byte, valLen)
		_, err = file.Read(valBytes)
		if err != nil {
			return nil, err
		}
		if key != nil {
			valBytes, err = utils.Decrypt(valBytes, key)
			if err != nil {
				return nil, err
			}
		}
		valDecoded, err := snappy.Decode(nil, valBytes)
		if err != nil {
			return nil, err
		}
		result[string(keyDecoded)] = string(valDecoded)
	}
	return result, nil
}
