// Package filewatcher
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package filewatcher

import "testing"

func TestWatchFiles(t *testing.T) {
	err := WatchFiles("../filehandler/testdir/testsrc", "../filehandler/testdir/testdes")
	if err != nil {
		t.Errorf("Error starting filewatcher.")
	}
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
