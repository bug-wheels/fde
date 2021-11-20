package util

import "os"

type FileEncode interface {
	encode(sourceFile, destinationFile string) error

	decode(sourceFile, destinationFile string) error
}

type FileHeadTitle struct {
	Algorithm string `json:"algorithm"`
	PKey      string `json:"pkey"`
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
