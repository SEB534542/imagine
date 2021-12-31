// Package imagine implements logic to encrypt and decrypt files as images.
package imagine

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var _ = log.Printf // TODO: remove when done

const (
	lengthMin int = 2000000 // minimum byte length
	lengthMax int = 5000000 // maximum byte length
)

var (
	prefix       string = "IMG_"   // prefix + counter + postfix are used in targetFname()
	counter             = 1000     // prefix + counter + postfix are used in targetFname()
	postfix      string = ".jpg"   // prefix + counter + postfix are used in targetFname()
	targetFolder string = "Photos" // target location for the generated files
)

type Key struct {
	origName string
	files    []string
}

// splitFile takes a file name, splits the file into chunks of bytes and returns
// a slice containing all chunks and an error.
func splitFile(name string) (output [][]byte, err error) {
	// Open file
	file, err1 := os.Open(name)
	if err1 != nil {
		err = fmt.Errorf("Error opening file '%v' in splitFile():\n%v", name, err1)
		return
	}
	defer file.Close()
	// Get file stat to determine total lenght of file
	fstat, _ := file.Stat()
	totalLen := int(fstat.Size())
	var off int // byte offset to read file the at offset
	for off < totalLen {
		chunk := random()
		if x := totalLen - off; chunk > x {
			chunk = x
		}
		byteBuff := make([]byte, chunk)
		bytesRead, err1 := file.ReadAt(byteBuff, int64(off))
		if err1 != nil {
			err = fmt.Errorf("Error reading file '%v' at %v splitFile():\n%v", name, off, err1)
			return
		}
		off += bytesRead
		output = append(output, byteBuff)
	}
	return
}

// random returns a random number for the []byte length. The int has a value
// between constant lenghtMin and lenghtMax.
func random() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(lengthMax-lengthMin+1) + lengthMin

}

// reverse takes a slice of byte and returns the same slice in reversed order.
func reverse(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

// encrypt takes a slice of byte, changes each byte to next byte (+1) and
// returns the updated slice of byte.
func encrypt(b []byte) []byte {
	for i, _ := range b {
		switch b[i] {
		case 255:
			b[i] = 0
		default:
			b[i] += 1
		}
	}
	return b
}

// decrypt takes a slice of byte, changes each byte to previous byte (-1) and
// returns the updated slice of byte.
func decrypt(b []byte) []byte {
	for i, _ := range b {
		switch b[i] {
		case 0:
			b[i] = 255
		default:
			b[i] -= 1
		}
	}
	return b
}

// storeImage takes a filename and stores the corresponding
// data in the created filename.
func storeImage(fname string, data []byte) error {
	// Create file
	file, err := os.Create(fname)
	defer file.Close()
	if err != nil {
		return err
	}
	// Write test data in file
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// targetFname returns a string representing a filename, based on const prefix +
// var counter + const postfix and adds a 1 to counter.
func targetFname() string {
	// TODO: check if targetFname exists, if so, add 1000(?)
	c := fmt.Sprint(counter)
	if len(c) == 4 {
		c = fmt.Sprint("0" + c)
	}
	c = fmt.Sprintf("%v%v%v", prefix, c, postfix)
	counter += 1
	return c
}

// imageFile takes a file name and a target directory, splits that filename
// into a target directory and returns the key.
func imageFile(fname, targetFolder string) ([]string, error) {
	key := []string{}
	output, err := splitFile(fname)
	if err != nil {
		return key, fmt.Errorf("Error splitting '%v':\n%v\n", fname, err)
	}
	for _, v := range output {
		v = encrypt(reverse(v))
		targetFname := targetFname()
		err = storeImage(targetFolder+"\\"+targetFname, v)
		if err != nil {
			return key, fmt.Errorf("Error storing file '%v' into '%v' at '%v':\n%v\n", fname, targetFname, targetFolder, err)
		}
		key = append(key, targetFname)
	}
	return key, nil
}

// Imagine takes a slice of directories and returns the key
// (ie which images belong to which file as part of which directory),
// e.g. key[dir1][file1][1.jpg, 2.jpg].
func Imagine(dirs []string) map[string]map[string][]string {
	key := make(map[string]map[string][]string)
	// create output folder
	newDir(targetFolder)
	for _, dir := range dirs {
		dirRel := lastSegment(dir)
		if key[dirRel] == nil {
			key[dirRel] = map[string][]string{}
		}
		fnames, err := getFnames(dir)
		if err != nil {
			log.Printf("ERROR! Directory '%v' skipped due to error:\n%v", dir, err)
		} else {
			for _, fname := range fnames {
				output, err := imageFile(fname, targetFolder)
				if err != nil {
					log.Printf("ERROR! File '%v' in dir %v not included due to error:\n%v", fname, dir, err)
				} else {
					// make fname and dir relative and add to key
					fname = relPath(dir, fname)
					key[dirRel][fname] = output
				}
			}
		}
	}
	return key
}
