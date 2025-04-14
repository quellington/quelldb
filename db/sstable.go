// sorted string table 
package db

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func WriteSSTable(path string, data map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for k, v := range data {
		_, err := file.WriteString(fmt.Sprintf("%s|%s\n", k, v))
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadSSTable(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "|", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result, nil
}
