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

// func TestSaveAndReadGobMap(t *testing.T) {
// 	fname := "test2.jpg"
// 	a := map[string]map[string][]string{}
// 	a["dir"] = map[string][]string{
// 		"file": []string{"1.jpg", "2.jpg"},
// 	}
// 	// Save var a
// 	saveToGob(a, fname)
// 	// load var a into
// 	var b map[string]map[string][]string
// 	readGob(&b, fname)
// 	t.Log(a)
// 	t.Log(b)
// 	err := os.Remove(fname)
// 	if err != nil {
// 		t.Errorf("Unable to remove %v", err)
// 	}
// }
