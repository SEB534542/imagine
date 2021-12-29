// Package imagine implements logic to encrypt and decrypt files as images.
package imagine

import (
	"fmt"
	"log"
	"os"
)

// TODO: remove when done
var (
	_ = fmt.Printf
	_ = log.Printf
)

// func to read identify the tree inside a directory

// isFile takes a file path and returns true if it is a file.
func isFile(fname string) bool {
	info, err := os.Stat(fname)
	if os.IsNotExist(err) {
		log.Fatal("Item does not exist:", err)
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

// isDir takes a directory path and returns true if it is directory.
func isDir(fname string) bool {
	info, _ := os.Stat(fname)
	if info.IsDir() {
		return true
	}
	return false
}

// newDir creates a new directory. E.g. "c:\Directory\Folder" will create "Folder" in "Directory".
// If "Folder" already exists, it will NOT return an error. However, if "Directory" does not exist, it will return an error.
func newDir(dir string) error {
	err := os.Mkdir(dir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
