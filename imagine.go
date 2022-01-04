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

const (
	lengthMin int = 2000000 // minimum byte length
	lengthMax int = 5000000 // maximum byte length
)

// Encrypted images are named as follows:  Prefix + Counter + Postfix,
// through targetFname.
var (
	Prefix  string = "IMG_" // E.g. "IMG_"
	Counter        = 1000   // Any number with a lenght less than 5 will be formatted with leading zeros. E.g. 500 will be formatted to 00500
	Postfix string = ".jpg" // E.g. ".jpg"
)

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
func imageFile(fname, trg string) ([]string, error) {
	key := []string{}
	output, err := splitFile(fname)
	if err != nil {
		return key, fmt.Errorf("Error splitting '%v':\n%v\n", fname, err)
	}
	for _, v := range output {
		v = transform(v)
		trgFname := targetFname()
		err = storeImage(trg+"\\"+trgFname, v)
		if err != nil {
			return key, fmt.Errorf("Error storing file '%v' into '%v' at '%v':\n%v\n", fname, trgFname, trg, err)
		}
		key = append(key, trgFname)
	}
	return key, nil
}

// Imagine takes a slice of directories and creates images in ImageFolder
// containing the transformed data from all the files in the directories.
// In case of errors, the file is skipped and added to the error message returned.
func Imagine(dirs []string, trg string) (err error) {
	key := make(map[string]map[string][]string)
	// create output folder
	newDir(trg)
	// get fname to store key (ie first image)
	keyFname := targetFname()
	for _, dir := range dirs {
		dirRel := lastSegment(dir)
		if key[dirRel] == nil {
			key[dirRel] = map[string][]string{}
		}
		fnames, err1 := getFnames(dir)
		if err1 != nil {
			err = fmt.Errorf("%vERROR! Directory '%v' skipped due to error:\n%v\n", err, dir, err1)
		} else {
			for _, fname := range fnames {
				output, err1 := imageFile(fname, trg)
				if err1 != nil {
					log.Printf("ERROR! File '%v' in dir %v not included due to error:\n%v", fname, dir, err1)
				} else {
					// make fname and dir relative and add to key
					fname = relPath(dir, fname)
					key[dirRel][fname] = output
				}
			}
		}
	}
	saveToGob(key, trg+"\\"+keyFname)
	return
}

// deImageFile takes a file name, source folder and list of all 'images' that
// make up the file. It transforms the image(s) into the file and returns any
// errors.
func deImageFile(fname, src string, imgFnames []string) error {
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
		imgFname = src + "\\" + imgFname
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

// DeImagine takes a source folder, target folder and a key to transform the
// filesin the source folder and create the files in the target folder using
// the key stored in the keyFname and returns an error.
func DeImagine(src, trg, keyFname string) (err error) {
	var key map[string]map[string][]string
	newDir(trg)
	// Read key
	readGob(&key, src+"\\"+keyFname)
	// For each directory
	for dir, fnames := range key {
		// Create new directory
		dir = trg + "\\" + dir
		err1 := newDir(dir)
		if err1 != nil {
			err = fmt.Errorf("%vError creating new dir '%v': %v\n", err, dir, err1)
			continue
		}
		for fname, imgFnames := range fnames {
			checkSubdirs(dir + "\\" + fname)
			err1 = deImageFile(dir+"\\"+fname, src, imgFnames)
			if err1 != nil {
				err = fmt.Errorf("%vError creating file '%v' in dir '%v': %v\n", err, fname, dir, err1)
			}
		}
	}
	return
}

// storeFile takes a filename and stores the corresponding
// data in the created filename.
func storeFile(fname string, data []byte) error {
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
