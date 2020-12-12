package utils

import (
	"errors"
	"os"
)

// CheckExists check if a file/directory exist
func CheckExists(path string) (bool, error) {
	// Check if path is empty
	if len(path) == 0 {
		return false, errors.New("Invalid Path or Filename")
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
