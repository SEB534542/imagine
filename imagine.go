// Package imagine implements logic to encrypt and decrypt files as images.
package imagine

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

// TODO v2: add a counter for progress tracking (how many files and how many done)

var _ = log.Printf // TODO: remove when done

const (
	lengthMin int = 2000000 // minimum byte length
	lengthMax int = 5000000 // maximum byte length
)

var (
	Prefix      string = "IMG_"   // prefix + counter + postfix are used in targetFname()
	Counter            = 1000     // prefix + counter + postfix are used in targetFname()
	Postfix     string = ".jpg"   // prefix + counter + postfix are used in targetFname()
	ImageFolder string = "Photos" // target location for the generated files
)

var outputFolder string // represents the folder where the output from deImagine() should be stored

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

// transform takes a slice of byte, and returns the reversed and encrypted slice.
func transform(b []byte) []byte {
	return encrypt(reverse(b))
}

// deTransform takes a slice of byte, and returns the de-encrypted and reversed
// slice.
func deTransform(b []byte) []byte {
	return reverse(decrypt(b))
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
	c := fmt.Sprint(Counter)
	if len(c) == 4 {
		c = fmt.Sprint("0" + c)
	}
	c = fmt.Sprintf("%v%v%v", Prefix, c, Postfix)
	Counter += 1
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
		v = transform(v)
		targetFname := targetFname()
		err = storeImage(targetFolder+"\\"+targetFname, v)
		if err != nil {
			return key, fmt.Errorf("Error storing file '%v' into '%v' at '%v':\n%v\n", fname, targetFname, targetFolder, err)
		}
		key = append(key, targetFname)
	}
	return key, nil
}

/* Imagine takes a slice of directories and creates images in ImageFolder
containing the data from all the files in the directories. In case of errors,
the file is skipped */
func Imagine(dirs []string) {
	key := make(map[string]map[string][]string)
	// create output folder
	newDir(ImageFolder)
	// get fname to store key (ie first image)
	keyFname := targetFname()
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
				output, err := imageFile(fname, ImageFolder)
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
	// TODO: Transform and store key as a the first image
	storekey(keyFname)

	return
}

// deImageFile takes ... and returns the ...
func deImageFile(fname string, imgFnames []string) error {
	// Create "original" file
	file, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("Error opening file '%v':\n%v\n", fname, err)
	}
	defer file.Close()
	// Read, convert and write each image into the "original" file
	var off int
	for _, imgFname := range imgFnames {
		// Read
		b, err := ioutil.ReadFile(imgFname)
		if err != nil {
			return fmt.Errorf("Error transforming file '%v' back into '%v'\n%v\n", imgFname, fname, err)
		}
		// Transform
		b = deTransform(b)
		// Write
		_, err = file.WriteAt(b, int64(off))
		if err != nil {
			return fmt.Errorf("Error writing file '%v' back into '%v'\n%v\n", imgFname, fname, err)
		}
		off += len(b)
	}
	return nil
}

// TODO: Description DeImagine
func DeImagine() {
	// TODO: read all files from ImageFolder

	// TODO: take first file and transform to key

	// TODO: for each file in the key: deImagine()

	return
}

// TODO: Transform and store key as a the first image
func storekey(keyFname string) error {
	return nil
}

// TODO: get key from first file in output folder
func getKey() error {
	return nil
}
