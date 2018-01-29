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

	info, err := os.Stat(srcfp)
	if err != nil {
		return err
	}
	fmt.Printf("Src file perms: %s.\n", info.Mode())

	from, err := os.Open(srcfp)
	if err != nil {
		return err
	}
	defer from.Close()

	to := &os.File{}
	if ok := pathExists(desfp); !ok {
		fmt.Printf("File '%s' doesn't exist, creating it now...\n", desfp)
		to, err = os.Create(desfp)
		if err != nil {
			return err
		}
	} else {
		to, err = os.OpenFile(desfp, os.O_RDWR|os.O_CREATE, info.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(to, from)
		if err != nil {
			return err
		}
		defer to.Close()
	}

	fmt.Println(from.Name(), "has been copied to", to.Name())

	return nil
}

// CopyDir copies the source directory to the destination directory
func CopyDir(srcdir string, desdir string) error {
	if ok := pathExists(desdir); ok {
		fmt.Println("Path already exists, so no need to create directories.")
		return nil
	}
	fmt.Println("Path doesn't exist, so creating path now...")

	info, err := os.Stat(srcdir)
	if err != nil {
		return err
	}

	var desdirs []string
	if info.IsDir() {
		desdirs = splitPath(desdir, false)
	} else {
		desdirs = splitPath(desdir, true)
	}
	srcdirs := splitPath(srcdir, true)

	var srcpath, despath string
	for i, desdir := range desdirs {
		if i == 0 {
			despath = strings.Join([]string{despath, desdir}, "")
			srcpath = strings.Join([]string{srcpath, srcdirs[i]}, "")
		} else {
			despath = strings.Join([]string{despath, desdir}, "/")
			if i < len(srcdirs) {
				fmt.Println("Length is good.")
				srcpath = strings.Join([]string{srcpath, srcdirs[i]}, "/")
				info, err = os.Stat(srcpath)
				if err != nil {
					return err
				}
			} else {
				fmt.Println("Length isn't good.")
			}
		}
		if ok := pathExists(despath); !ok {
			fmt.Println("Creating Path:", despath)
			fmt.Println("Path mode:", info.Mode())
			if err := os.Mkdir(despath, info.Mode()); err != nil {
				return err
			}
		}
	}

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
