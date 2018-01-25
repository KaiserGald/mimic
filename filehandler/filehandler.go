// Package filehandler
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package filehandler

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// CopyFile will copy the supplied file to the supplied destination
func CopyFile(srcfp, desfp string) error {

	err := CopyDir(srcfp, desfp)
	if err != nil {
		return err
	}

	from, err := os.Open(srcfp)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(desfp, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	fmt.Println(from.Name(), "has been copied to", to.Name())

	return nil
}

// CopyDir copies the source directory to the destination directory
func CopyDir(srcdir string, desdir string) error {
	if ok := pathExists(desdir); ok {
		fmt.Println("exists")
		return nil
	}

	info, err := os.Stat(srcdir)
	if err != nil {
		return err
	}

	var dirs []string
	if info.IsDir() {
		dirs = splitPath(desdir, false)
	} else {
		dirs = splitPath(desdir, true)
	}

	var path string
	for i, dir := range dirs {
		if i == 0 {
			path = strings.Join([]string{path, dir}, "")
		} else {
			path = strings.Join([]string{path, dir}, "/")
		}
		if ok := pathExists(path); !ok {
			fmt.Println(ok)
			if err := os.Mkdir(path, info.Mode()); err != nil {
				return err
			}
		}
	}
	fmt.Println(path)

	return nil
}

// Remove removes the given file or directory
func Remove(fp string) error {
	if err := os.Remove(fp); err != nil {
		return err
	}
	return nil
}

// Rename renames the file or directory to the given name
func Rename(old, new string) error {
	err := os.Rename(old, new)
	if err != nil {
		return err
	}
	return nil
}

// Chmod changes the given file's permissions to the value passed in
func Chmod(src, des string) error {
	f1, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.Chmod(des, f1.Mode()); err != nil {
		return err
	}

	return nil
}

// pathExists checks to see if the given path exists, returns a bool
func pathExists(fp string) bool {
	_, err := os.Stat(fp)
	if err != nil {
		return false
	}
	return true
}

// splitPath splits the path string into the individual files and returns them in an array
func splitPath(path string, dirs bool) []string {
	sp := strings.Split(path, "/")
	if dirs {
		sp = sp[:len(sp)-1]
	}

	return sp
}
