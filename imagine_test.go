package imagine

import (
	"fmt"
	"os"
	"testing"
)

func TestRandom(t *testing.T) {
	for i := 0; i < 1000; i++ {
		x := random()
		if x < lengthMin || x > lengthMax {
			t.Errorf("Want value within %v and %v. Got: %v", lengthMin, lengthMax, x)
		}
	}
}

func TestSplitFile(t *testing.T) {
	want := "sfsdgdfhfghdfg ddf gdgdf gdfgfgffg3534"
	fnameWant, fnameGot := "test_want.txt", "test_got.txt"
	// Create file for testing
	file, err := os.Create(fnameWant)
	if err != nil {
		t.Errorf("Error creating file '%v':\n%v", fnameWant, err)
	}
	// Write test data in file
	_, err = file.Write([]byte(want))
	file.Close()
	// Split file
	output, err := splitFile(fnameWant)
	if err != nil {
		t.Errorf("Error splitting file '%v':\n'%v'", fnameWant, err)
	}
	// Evaluate want with got
	var got []byte
	for _, v := range output {
		got = append(got, v...)
	}
	if string(got) != want {
		// Creating file to compare want with got offline
		file, err = os.Create(fnameGot)
		defer file.Close()
		if err != nil {
			t.Errorf("Error creating file '%v':\n%v", fnameGot, err)
		}
		_, err = file.Write(got)
		t.Errorf("Error while comparing test file with output, please compare %v and %v", fnameGot, fnameWant)
	} else {
		// No error, remove created file
		os.Remove(fnameWant)
	}
}

func TestReverse(t *testing.T) {
	cases := []struct{ input, want []byte }{
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{[]byte{0, 1, 2, 3, 4}, []byte{4, 3, 2, 1, 0}},
		{[]byte{22, 28, 100, 5, 105}, []byte{105, 5, 100, 28, 22}},
		{[]byte{22, 28, 5, 105}, []byte{105, 5, 28, 22}},
	}
	for x, c := range cases {
		got := reverse(c.input)
		errMsg := fmt.Sprintf("Error reversing case %v.\nWant: %v\nGot: %v", x, c.want, got)
		if len(c.want) != len(got) {
			t.Errorf(errMsg)
		}
		for i, _ := range c.want {
			if c.want[i] != got[i] {
				t.Errorf(errMsg)
				break
			}
		}
	}
}

func TestEncrypt(t *testing.T) {
	cases := []struct{ input, want []byte }{
		{[]byte{0, 1, 2, 3, 4, 5, 6}, []byte{1, 2, 3, 4, 5, 6, 7}},
		{[]byte{255, 1, 2, 3, 4, 5, 6}, []byte{0, 2, 3, 4, 5, 6, 7}},
	}
	for i, c := range cases {
		got := encrypt(c.input)
		for j, _ := range c.want {
			if c.want[j] != got[j] {
				t.Errorf("Error in encrypting case %v.\nWant: %v\nGot: %v\n", i, c.want, got)
				break
			}
		}
	}
}

func TestDecrypt(t *testing.T) {
	cases := []struct{ input, want []byte }{
		{[]byte{0, 1, 2, 3, 4, 5, 6}, []byte{255, 0, 1, 2, 3, 4, 5}},
		{[]byte{255, 1, 2, 3, 4, 5, 6}, []byte{254, 0, 1, 2, 3, 4, 5}},
	}
	for i, c := range cases {
		got := decrypt(c.input)
		for j, _ := range c.want {
			if c.want[j] != got[j] {
				t.Errorf("Error in decrypting case %v.\nWant: %v\nGot: %v\n", i, c.want, got)
				break
			}
		}
	}
}

func TestEncryptDecrypt(t *testing.T) {
	cases := []struct{ input, want []byte }{
		{[]byte{0, 1, 2, 3, 4, 5, 6}, []byte{0, 1, 2, 3, 4, 5, 6}},
		{[]byte{255, 1, 2, 3, 4, 5, 6}, []byte{255, 1, 2, 3, 4, 5, 6}},
	}
	for i, c := range cases {
		got := decrypt(encrypt(c.input))
		for j, _ := range c.want {
			if c.want[j] != got[j] {
				t.Errorf("Error in encrypting + decrypting case %v.\nWant: %v\nGot: %v\n", i, c.want, got)
				break
			}
		}
	}
}

func TestStoreFile(t *testing.T) {
	//fname, target string, data []byte
	fname := "test.txt"
	target := "."
	err := storeImage(target+"\\"+fname, []byte("test"))
	if err != nil {
		t.Errorf("Error storing file '%v'", fname)
	}
	os.Remove(target + "\\" + fname)
}

func TestTargetFile(t *testing.T) {
	Counter = 1000
	cases := []struct {
		counter int
		want1   string
		want2   string
	}{
		{1000, "IMG_01000.jpg", "IMG_01001.jpg"},
		{9999, "IMG_09999.jpg", "IMG_10000.jpg"},
		{99999, "IMG_99999.jpg", "IMG_100000.jpg"},
	}
	for _, c := range cases {
		Counter = c.counter
		// Round 1
		got := targetFname()
		if got != c.want1 {
			t.Errorf("Target filename incorrect. Want: '%v'. Got: '%v'", c.want1, got)
		}
		// Round 2
		got = targetFname()
		if got != c.want2 {
			t.Errorf("Target filename incorrect. Want: '%v'. Got: '%v'", c.want2, got)
		}
	}
}

func ExampleImagine() {
	fmt.Println(imagine.Imagine([]string{".\\Test files"}, ".\\Photos"))
}

func ExampleDeImagine() {
	fmt.Println(imagine.DeImagine(".\\Photos", ".\\Output", "IMG_01000.jpg"))
}
