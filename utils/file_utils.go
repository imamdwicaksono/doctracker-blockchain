package utils

import (
	"io/ioutil"
	"os"
)

func CreateDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteToFile writes the given content to the specified filename.
func WriteToFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), os.ModePerm)
}
