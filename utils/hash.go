package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
)

// FileCategoryHash returns SHA-256 of file data concatenated with category.
func FileCategoryHash(path, category string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(append(data, []byte(category)...))
	return hex.EncodeToString(h[:]), nil
}
