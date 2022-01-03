package imagine

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

// saveGob encodes an interface and stores it as a Gob into a file named fname.
func saveToGob(i interface{}, fname string) error {
	var data bytes.Buffer

	enc := gob.NewEncoder(&data) // Will write to data
	//	dec := gob.NewDecoder(&data) // Will read from data

	// Encode (send) some values.
	err := enc.Encode(i)
	if err != nil {
		return fmt.Errorf("Error encoding '%v': %v", fname, err)
	}

	// Tranform data
	b := transform(data.Bytes())

	// Store data
	err = ioutil.WriteFile(fname, b, 0644)
	if err != nil {
		return fmt.Errorf("Error storing '%v': %v", fname, err)
	}
	return nil
}

// readGob reads a gob from a file and converts it into an interface. It takes a pointer to an interface and a file name.
func readGob(i interface{}, fname string) error {
	// Initialize decoder
	var data bytes.Buffer
	dec := gob.NewDecoder(&data) // Will decode (read) and store into data

	// Read content from file
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("Error reading file '%v': %v", fname, err)
	}

	// Transform data
	b := deTransform(content)

	y := bytes.NewBuffer(b)
	data = *y

	// Decode (receive) and print the values.

	err = dec.Decode(i)
	if err != nil {
		return fmt.Errorf("Error decoding into '%v': %v (%v)", fname, err, i)
	}
	return nil
}
