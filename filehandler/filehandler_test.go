// Package filehandler_test
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package filehandler

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestCopyFile(t *testing.T) {
	src := "testdir/testsrc/test.txt"
	des := "testdir/testdes/testcreate.txt"
	err := CopyFile(src, des)
	if err != nil {
		t.Errorf("Error copying '%s' to '%s': %v\n", src, des, err)
	}

	if ok := exists(des); !ok {
		t.Errorf("File didn't copy over to destination.\n")
	}

	os.Remove(des)

	des = "testdir/testdes/test.txt"
	os.Create(des)
	err = CopyFile(src, des)
	if err != nil {
		t.Errorf("Error copying '%s' to '%s': %v\n", src, des, err)
	}

	f1, _ := ioutil.ReadFile(src)
	f2, _ := ioutil.ReadFile(des)

	if bytes.Compare(f1, f2) != 0 {
		t.Errorf("Error copying: '%s' and '%s' are not identical.\n", src, des)
	}
	os.Remove(des)
}

func TestCopyDir(t *testing.T) {
	src := "testdir/testsrc"
	des := "testdir/testdircopy"
	err := CopyDir(src, des)
	if err != nil {
		t.Errorf("Error copying '%s' to '%s': %v\n", src, des, err)
	}

	if ok := exists(des); !ok {
		t.Errorf("Directory didn't copy over to destination.\n")
	}

	os.Remove(des)
}

func TestRemove(t *testing.T) {
	file := "test.txt"
	os.Create(file)
	err := Remove(file)
	if err != nil {
		t.Errorf("Error removing '%s': %v\n", file, err)
	}

	if ok := exists(file); ok {
		t.Errorf("'%s' still exists.\n", file)
		os.Remove(file)
	}
}

func TestRename(t *testing.T) {
	old := "test.txt"
	new := "test1.txt"
	os.Create(old)
	err := Rename(old, new)
	if err != nil {
		t.Errorf("Error renaming '%s' to '%s': %v\n", old, new, err)
	}

	if ok := exists(new); !ok {
		t.Errorf("File was not renamed: %v\n", err)
	}
	os.Remove(old)
	os.Remove(new)

	old = "renamedir"
	new = "renamedirtest"
	os.Mkdir(old, 0700)
	err = Rename(old, new)
	if err != nil {
		t.Errorf("Error renaming '%s' to '%s': %v\n", old, new, err)
	}

	if ok := exists(new); !ok {
		t.Errorf("Directory was not renamed: %v\n", err)
	}
	os.Remove(old)
	os.Remove(new)

}

func TestChmod(t *testing.T) {
	src := "test.txt"
	des := "testdes.txt"
	var perm os.FileMode = 0777
	os.Create(src)
	os.Create(des)
	os.Chmod(src, perm)

	err := Chmod(src, des)
	if err != nil {
		t.Errorf("Error changing file mode: %v\n", err)
	}

	f1, _ := os.Stat(src)
	f2, _ := os.Stat(des)

	if f1.Mode() != f2.Mode() {
		t.Errorf("Permissions between the files are different.\n")
	}

	os.Remove(src)
	os.Remove(des)
}

func exists(fp string) bool {
	_, err := os.Stat(fp)
	if err != nil {
		return false
	}
	return true
}
