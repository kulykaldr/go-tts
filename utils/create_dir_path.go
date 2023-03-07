package utils

import (
	"log"
	"os"
	"path/filepath"
)

func CreateDirPath(dirPath ...string) string {
	p := filepath.Join(dirPath...)
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal("error create dir", err)
	}
	dir, _ := filepath.Abs(p)

	return dir
}
