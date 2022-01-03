package imagine

import (
	"os"
	"testing"
)

func TestSaveAndReadGob(t *testing.T) {
	type P struct {
		X, Y, Z int
		Name    string
	}
	fname := "test.jpg"
	a := P{3, 4, 5, "Pythagoras"}
	// Save var a (which is of type P)
	saveToGob(a, fname)

	// load var a into struct P
	var b P
	readGob(&b, fname)
	if !(a.X == b.X && a.Y == b.Y && a.Z == b.Z && a.Name == b.Name) {
		t.Error("Data not correctly saved and/or loaded")
	}
	err := os.Remove(fname)
	if err != nil {
		t.Errorf("Unable to remove %v", err)
	}
}
