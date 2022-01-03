// Package imagine implements logic to encrypt and decrypt files as images.
package imagine

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// newDir takes a folder path, creates a new directory and returns an error if
// it cannot create the new file and does not already exist.
func newDir(dir string) error {
	err := os.Mkdir(dir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// isFile takes a file path and returns true if it represents a file.
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

// isDir takes a directory path and returns true if it represents a directory.
func isDir(fname string) bool {
	info, _ := os.Stat(fname)
	if info.IsDir() {
		return true
	}
	return false
}

// getFnames retrieves the subdirectories and files within the source folder and stores it in Subdirs and Files.
func getFnames(dir string) (fnames []string, err error) {
	// Open directory
	file, err1 := os.Open(dir)
	if err1 != nil {
		err = fmt.Errorf("Unable to open '%v': %v", dir, err1)
		return
	}
	defer file.Close()

	// Read all files and directories
	list, _ := file.Readdirnames(0)
	for _, item := range list {
		if itemPath := dir + "\\" + item; isFile(itemPath) {
			// File
			fnames = append(fnames, itemPath)
		} else {
			// Directory
			subFnames, err1 := getFnames(itemPath)
			if err != nil {
				err = fmt.Errorf("Unable to getFnames from '%v': %v", dir, err1)
				return
			}
			fnames = append(fnames, subFnames...)
		}
	}
	return
}

// lastSegment takes a path and returns the last segment of that path.
// E.g. "c:\test" will return "test".
func lastSegment(p string) string {
	v := strings.Split(p, "\\")
	if n := v[len(v)-1]; n != "" {
		return n
	}
	return v[len(v)-2]
}

// relPath takes two paths and returns the path without the root.
func relPath(root, path string) string {
	// Check and add backslash to root for comparison to path
	if string(root[len(root)-1]) != "\\" {
		root += "\\"
	}
	// Check that path is in root
	if root != path[:len(root)] {
		log.Fatalf("Path not in root.\nRoot:\t%v\nPath:\t%v\n", root, path)
	}
	//	Remove root from path
	path = path[len(root):]
	return path
}

// checkSubdirs checks every level of the fname path if the subdirectory exists.
func checkSubdirs(fname string) {
	// TODO: check if there is a more efficient way then re-checking every time
	// Split fname path
	x := strings.Split(fname, "\\")
	if len(x) < 2 {
		return
	}
	path := x[0]
	for _, v := range x[1 : len(x)-1] {
		path += "\\" + v
		newDir(path)
	}
}
