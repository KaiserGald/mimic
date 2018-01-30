// Package filewatcher
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package filewatcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KaiserGald/mimic/filehandler"
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

func TestInitializeFileTree(t *testing.T) {
	os.MkdirAll("testsrc/subtest/subtest1", 0777)
	os.MkdirAll("testsrc/subtest/subtest2", 0777)
	os.Create("testsrc/test.txt")
	os.Create("testsrc/test1.txt")
	os.Create("testsrc/subtest/test.txt")
	os.Create("testsrc/subtest/subtest1/test.txt")
	os.Create("testsrc/subtest/subtest1/test1.txt")
	os.Create("testsrc/subtest/subtest2/test.txt")
	os.Create("testsrc/subtest/subtest2/test1.txt")

	err := initializeFileTree(srcfp, desfp, relfp)
	if err != nil {
		t.Errorf("Error initializing file tree: %v", err)
	}

	expected, _ := mapTree("testsrc")
	fmt.Println("Expected:", expected["subtest/test.txt"].Name())

	actual, _ := mapTree("testdes")

	if len(actual) == 0 {
		t.Errorf("Copied file tree was not mapped")
	}

	compareTrees(expected, actual, t)

	os.RemoveAll("testsrc")
	os.Mkdir("testsrc", 0770)
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

	os.Remove(srcfp + filename)
	os.RemoveAll(testdir)
	os.Remove(desfp + filename)
	os.RemoveAll(desfp + "/testdir")

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
	oldname := "/test.txt"
	newname := "/rename.txt"
	des := desfp + oldname
	newdes := desfp + newname
	oldsrc := srcfp + oldname
	newsrc := srcfp + newname
	os.Create(oldsrc)
	os.Create(newsrc)
	os.Create(des)

	info, _ := os.Stat(oldsrc)
	os.Rename(oldsrc, newsrc)

	event := watcher.Event{
		watcher.Rename,
		relfp + oldname + " -> " + relfp + newname,
		info,
	}

	err := handleRename(event, desfp, relfp)
	if err != nil {
		t.Errorf("Error renaming file: %v\n", err)
	}

	info, _ = os.Stat(newdes)
	if "/"+info.Name() != newname {
		t.Errorf("File names don't match.\n")
	}

	os.Remove(newdes)
	os.Remove(newsrc)

}

func TestHandleChmod(t *testing.T) {
	filename := "/test.txt"
	srcpath := srcfp + filename
	despath := desfp + filename

	os.Create(srcpath)
	os.Create(despath)

	os.Chmod(srcpath, 0777)
	srcinfo, _ := os.Stat(srcpath)

	event := watcher.Event{
		watcher.Chmod,
		relfp + filename,
		srcinfo,
	}

	err := handleChmod(event, srcfp, desfp, relfp)
	if err != nil {
		t.Errorf("Error changing file permissions.")
	}

	desinfo, _ := os.Stat(despath)
	if desinfo.Mode() != srcinfo.Mode() {
		t.Errorf("File permissions do not match.\n")
	}

	dir := "/testdir"
	srcdir := srcfp + dir
	desdir := desfp + dir
	os.Mkdir(srcdir, 0770)
	os.Mkdir(desdir, 0660)

	srcinfo, _ = os.Stat(srcdir)
	event = watcher.Event{
		watcher.Chmod,
		relfp + dir,
		srcinfo,
	}

	err = handleChmod(event, srcfp, desfp, relfp)
	if err != nil {
		t.Errorf("Error changing directory permissions.")
	}

	desinfo, _ = os.Stat(desdir)

	if srcinfo.Mode() != desinfo.Mode() {
		t.Errorf("Directory permissions do not match.\n")
	}

	os.Remove(srcdir)
	os.Remove(desdir)
}

func TestHandleMove(t *testing.T) {
	filename := "/test.txt"
	desdir := "/test"
	srcf := srcfp + filename
	desf := srcfp + desdir + filename
	copysrcf := desfp + filename
	copydesf := desfp + desdir + filename
	fmt.Println(srcf)
	fmt.Println(desf)
	fmt.Println(copysrcf)
	fmt.Println(copydesf)

	os.Create(srcf)
	os.Create(desf)
	os.Mkdir(srcfp+desdir, 0777)
	os.Mkdir(desfp+desdir, 0777)

	filehandler.CopyFile(srcf, desf)
	info, _ := os.Stat(srcf)
	os.Remove(srcf)
	event := watcher.Event{
		watcher.Move,
		relfp + filename + " -> " + relfp + desdir + filename,
		info,
	}

	err := handleMove(event, srcfp, desfp, relfp)
	if err != nil {
		t.Errorf("Error moving file: %v\n", err)
	}

	_, err = os.Open(copydesf)
	if err != nil {
		t.Errorf("File was not moved\n")
	}

	f1, _ := ioutil.ReadFile(copysrcf)
	f2, _ := ioutil.ReadFile(copydesf)

	if bytes.Compare(f1, f2) != 0 {
		t.Errorf("Error moving: '%s' and '%s' are not identical.\n", copysrcf, copydesf)
	}

	os.RemoveAll(srcfp + desdir)
	os.Remove(copysrcf)
	os.RemoveAll(desfp + desdir)
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

func TestMapTree(t *testing.T) {

	os.MkdirAll("testsrc/subtest", 0777)
	os.Create("testsrc/test.txt")
	os.Create("testsrc/test1.txt")
	os.Create("testsrc/subtest/test.txt")
	ioutil.WriteFile("testsrc/test.txt", []byte("This is the string"), 0777)
	ioutil.WriteFile("testsrc/test1.txt", []byte("This is a different string"), 0777)
	ioutil.WriteFile("testsrc/subtest/test.txt", []byte("This is also a different one"), 0777)

	expected := make(map[string]os.FileInfo)

	filepath.Walk("testsrc", func(path string, info os.FileInfo, err error) error {
		files := strings.Split(path, "/")
		files = files[1:]
		path = strings.Join(files, "/")
		if len(path) != 0 {
			expected[path] = info
		}
		return nil
	})

	actual, err := mapTree("testsrc")
	if err != nil {
		t.Errorf("Error mapping file tree: %v", err)
	}
	if len(actual) == 0 {
		t.Errorf("Tree map was not created")
	}

	compareTrees(expected, actual, t)
}

func compareTrees(expected, actual map[string]os.FileInfo, t *testing.T) {
	for path := range actual {
		if actual[path].Name() != expected[path].Name() {
			t.Errorf("Names don't match, expected %v, got %v", expected[path].Name(), actual[path].Name())
		}
		if actual[path].Size() != expected[path].Size() {
			t.Errorf("File sizes don't match, expected %v, got %v", expected[path].Size(), actual[path].Size())
		}
		if actual[path].Mode() != expected[path].Mode() {
			t.Errorf("File permissions don't match, expected %v, got %v", expected[path].Mode(), actual[path].Mode())
		}
	}
}
