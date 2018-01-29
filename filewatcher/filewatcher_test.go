// Package filewatcher
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package filewatcher

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/radovskyb/watcher"
)

var (
	srcfp string
	desfp string
	relfp string
)

func TestMain(m *testing.M) {
	// set source and destination dirs and make them
	srcfp = "testsrc"
	desfp = "testdes"
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	relfp = strings.Join([]string{dir, "../" + srcfp}, "/")
	os.Mkdir(srcfp, 0770)
	os.Mkdir(desfp, 0770)

	r := m.Run()

	// clean up
	os.RemoveAll(srcfp)
	os.RemoveAll(desfp)

	os.Exit(r)
}

func TestInitWatcher(t *testing.T) {
	w, fp, err := initWatcher("testdir/testsrc")
	if w == nil {
		t.Errorf("Error creating watcher.")
	}

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	relfp := strings.Join([]string{dir, "testdir/testsrc"}, "/")
	if fp != relfp {
		t.Errorf("Filepaths do not match.")
	}

	if err != nil {
		t.Errorf("Error initializing watcher: %s", err)
	}
}

func TestHandleCreate(t *testing.T) {
	testdir := srcfp + "/test"
	testfile := testdir + "/test.txt"

	os.Mkdir(testdir, 0770)
	info, _ := os.Stat(testdir)
	event := watcher.Event{
		watcher.Create,
		relfp + "/test",
		info,
	}

	err := handleCreate(event, srcfp, desfp, relfp)
	if err != nil {
		t.Errorf("Error creating directory.")
	}

	_, err = os.Stat(desfp + "/test")
	if err != nil {
		t.Errorf("Directory was not copied: %s", err)
	}

	os.Create(testfile)

	info, _ = os.Stat(testfile)
	event = watcher.Event{
		watcher.Create,
		relfp + "/test/test.txt",
		info,
	}

	err = handleCreate(event, srcfp, desfp, relfp)
	if err != nil {
		t.Errorf("Error creating file.")
	}

	_, err = os.Stat(desfp + "/test/test.txt")
	if err != nil {
		t.Errorf("File was not copied: %s", err)
	}

	os.RemoveAll(testdir)
	os.RemoveAll(desfp + "/test")
}

func TestHandleWrite(t *testing.T) {
	filename := "/test.txt"
	testfile := srcfp + filename
	fulltestpath := relfp + filename
	os.Create(testfile)

	info, _ := os.Stat(testfile)
	event := watcher.Event{
		watcher.Write,
		fulltestpath,
		info,
	}

	err := handleWrite(event, srcfp, desfp, relfp)
	if err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	_, err = os.Stat(desfp + filename)
	if err != nil {
		t.Errorf("File was not copied: %v", err)
	}

	testdir := srcfp + "/testdir"
	testfile = testdir + filename
	os.MkdirAll(testdir, 0770)

	os.Create(testfile)
	fulltestpath = relfp + "/testdir" + filename

	info, _ = os.Stat(testfile)
	event = watcher.Event{
		watcher.Write,
		fulltestpath,
		info,
	}

	err = handleWrite(event, srcfp, desfp, relfp)
	if err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	//os.Remove(srcfp + filename)
	//os.RemoveAll(testdir)
	//os.Remove(desfp + filename)
	//os.RemoveAll(desfp + "/testdir")

}

func TestHandleRemove(t *testing.T) {
	filename := "/test.txt"
	srcpath := srcfp + filename
	despath := desfp + filename
	os.Create(srcpath)
	os.Create(despath)

	info, _ := os.Stat(srcpath)
	event := watcher.Event{
		watcher.Remove,
		relfp + filename,
		info,
	}

	os.Remove(srcpath)

	err := handleRemove(event, srcfp, desfp, relfp)
	if err != nil {
		t.Errorf("Error handling removal: %v", err)
	}

	_, err = os.Stat(despath)
	if err == nil {
		t.Errorf("File wasn't removed.")
	}

}

func TestHandleRename(t *testing.T) {

}

func TestHandleChmod(t *testing.T) {

}

func TestHandleMove(t *testing.T) {

}

func TestBuildPaths(t *testing.T) {
	ep := "/home/workspace/projectroot/test/testsrc/dir/test.txt"
	srcfp := "test/testsrc"
	desfp := "test/testdes"
	relfp := "/home/workspace/projectroot/test/testsrc"
	srctest := "test/testsrc/dir/test.txt"
	destest := "test/testdes/dir/test.txt"

	src, des := buildPaths(ep, srcfp, desfp, relfp)
	if src != srctest {
		t.Errorf("Error building source file path. Expected '%s' got '%s'.\n", srctest, src)
	}

	if des != destest {
		t.Errorf("Error building destination path. Expected '%s' got '%s'.\n", destest, des)
	}
}
