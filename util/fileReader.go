package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetDevHtml() ([]byte, error) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	testFile := exPath + "/test.txt"
	return ioutil.ReadFile(testFile)
}
