package imagine

import (
	"fmt"
	"log"

	//	"os"
	"testing"
)

// TODO: remove when done
var (
	_ = fmt.Printf
	_ = log.Printf
)

// func TestDir(t *testing.T) {
// 	dir := ".\\test_folder"
// 	fname := "temp.txt"
// 	// Create dir for testing
// 	err := newDir(dir)
// 	if err != nil {
// 		t.Errorf("Error creating new folder: '%v'\n%v", dir, err)
// 	}

// 	// Test if folder is a folder
// 	if b := isDir(dir); b != true {
// 		t.Errorf("Error check isDir: Want:'%v' got: '%v'", true, b)
// 	}

// 	// Create folder for testing
// 	file, err := os.Create(dir + "\\" + fname)
// 	file.Close()
// 	if err != nil {
// 		t.Errorf("Error creating file '%v':\n%v", fname, err)
// 	}

// 	// Test if file is a file
// 	if b := isFile(dir + "\\" + fname); b != true {
// 		t.Errorf("Error check isFile: want '%v' got '%v'", true, b)
// 	}

// 	//
// 	output, err := getFnames(dir)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if output[0] != dir+"\\"+fname {
// 		t.Errorf("Want: %v. Got: %v", dir+"\\"+fname, output[0])
// 	}

// 	// Remove file + directory
// 	os.Remove(dir + "\\" + fname)
// 	os.Remove(dir)
// }

func TestLastSegment(t *testing.T) {
	cases := []struct{ input, want string }{
		{"c:\\test", "test"},
		{"c:\\test ", "test "},
		{"c:\\test1\\test2\\test3", "test3"},
	}
	for _, c := range cases {
		got := lastSegment(c.input)
		if c.want != got {
			t.Errorf("Want: '%v'. Got: '%v'.", c.want, got)
		}
	}
}

func TestRelPath(t *testing.T) {
	cases := []struct{ root, path, want string }{
		{"c:\\test1\\test2", "c:\\test1\\test2\\test3", "test3"},
		{"c:\\test1\\test2\\", "c:\\test1\\test2\\test3", "test3"},
		{"c:\\test1", "c:\\test1\\test2\\test3", "test2\\test3"},
		{"c:\\test1\\test2\\test3", "c:\\test1\\test2\\test3\\test4", "test4"},
		{"c:\\test1", "c:\\test1\\test2\\test3\\test4", "test2\\test3\\test4"},
	}
	for _, c := range cases {
		got, err := relPath(c.root, c.path)
		if err != nil {
			t.Error(err)
		}
		if c.want != got {
			t.Errorf("Want: '%v' Got: '%v'", c.want, got)
		}
	}
}
