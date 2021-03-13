package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetLocalHtml(p string) ([]byte, error) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	var testFile string
	if p == "" {
		testFile = exPath + "/test.txt"
	} else {
		testFile = exPath + "/" + p
	}
	return ioutil.ReadFile(testFile)
}

func GetFilenamesFromDirectory(p string) []string {
	dir := filepath.Dir(p)
	if dir == "" {
		fmt.Printf("directory %s not found\n", p)
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(p)
	if err != nil {
		panic(err)
	}

	result := make([]string, 0)
	for _, f := range files {
		result = append(result, f.Name())
	}
	return result
}
