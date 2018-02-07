// Package filehandler
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package filehandler

import (
	"io"
	"os"
	"strings"

	"github.com/KaiserGald/logger"
)

var l *logger.Logger

// Init initializes the filehandler.
func Init(lg *logger.Logger) {
	l = lg
}

// CopyFile will copy the supplied file to the supplied destination
func CopyFile(srcfp, desfp string) error {

	l.Debug.Log("Copying directory '%v' to '%v'.", srcfp, desfp)
	err := CopyDir(srcfp, desfp)
	if err != nil {
		return err
	}
	l.Debug.Log("Done copying directory.")

	l.Debug.Log("Getting source file info.")
	info, err := os.Stat(srcfp)
	if err != nil {
		return err
	}
	l.Debug.Log("Done.")

	l.Debug.Log("Opening '%v'.", srcfp)
	from, err := os.Open(srcfp)
	if err != nil {
		return err
	}
	l.Debug.Log("'%v' successfully opened!", srcfp)
	defer from.Close()

	l.Debug.Log("Begin file copy...")
	to := &os.File{}
	l.Debug.Log("Checking if '%v' exists...", desfp)
	if ok := pathExists(desfp); !ok {
		l.Notice.Log("File '%s' doesn't exist, creating it now...", desfp)
		to, err = os.Create(desfp)
		if err != nil {
			return err
		}
		l.Debug.Log("File '%v' successfully created.", desfp)
	} else {
		l.Debug.Log("File already exists.")
		l.Debug.Log("Opening file '%v'.", desfp)
		to, err = os.OpenFile(desfp, os.O_RDWR|os.O_CREATE, info.Mode())
		if err != nil {
			return err
		}
		l.Debug.Log("File '%v' successfully opened!", desfp)
		l.Debug.Log("Copying file '%v' to '%v'.", srcfp, desfp)
		_, err = io.Copy(to, from)
		if err != nil {
			return err
		}
		l.Debug.Log("File successfully copied.")
		defer to.Close()
	}
	l.Debug.Log("File copy done.")

	return nil
}

// CopyDir copies the source directory to the destination directory
func CopyDir(srcdir, desdir string) error {
	l.Debug.Log("Does '%v' already exist?", desdir)
	if ok := pathExists(desdir); ok {
		l.Debug.Log("Yes!")
		l.Debug.Log("Directory already exists, so no need to create it.")
		return nil
	}
	l.Debug.Log("No!")
	l.Debug.Log("Directory doesn't exist, so creating it now...")

	l.Debug.Log("Getting file info for '%v'", srcdir)
	info, err := os.Stat(srcdir)
	if err != nil {
		return err
	}
	l.Debug.Log("Done.")

	l.Debug.Log("Begin copying directories...")
	l.Debug.Log("Building paths...")
	var desdirs []string
	if info.IsDir() {
		desdirs = splitPath(desdir, false)
	} else {
		desdirs = splitPath(desdir, true)
	}
	srcdirs := splitPath(srcdir, true)
	l.Debug.Log("srcdirs: %v", srcdirs)
	l.Debug.Log("desdirs: %v", desdirs)
	l.Debug.Log("Joining paths into a file path...")
	var srcpath, despath string
	for i, desdir := range desdirs {
		if i == 0 {
			despath = strings.Join([]string{despath, desdir}, "")
			srcpath = strings.Join([]string{srcpath, srcdirs[i]}, "")
		} else {
			despath = strings.Join([]string{despath, desdir}, "/")
			if i < len(srcdirs) {
				srcpath = strings.Join([]string{srcpath, srcdirs[i]}, "/")
				info, err = os.Stat(srcpath)
				if err != nil {
					return err
				}
			}
		}
		l.Debug.Log("Does directory '%v' exist?", despath)
		if ok := pathExists(despath); !ok {
			l.Notice.Log("Directory '%v' doesn't exist, creating it now...", despath)
			if err := os.Mkdir(despath, info.Mode()); err != nil {
				return err
			}
		}
	}
	l.Debug.Log("Done copying directories.")

	return nil
}

// Remove removes the given file or directory
func Remove(fp string) error {
	l.Debug.Log("Removing '%v' now...", fp)
	if err := os.Remove(fp); err != nil {
		return err
	}
	l.Debug.Log("File '%v' successfully removed!", fp)
	return nil
}

// Rename renames the file or directory to the given name
func Rename(old, new string) error {
	l.Debug.Log("Renaming '%v' now...", old)
	err := os.Rename(old, new)
	if err != nil {
		return err
	}
	l.Debug.Log("Done renaming. '%v' is now '%v'.", old, new)
	return nil
}

// Chmod changes the given file's permissions to the value passed in
func Chmod(src, des string) error {
	l.Debug.Log("Getting file info...")
	f1, err := os.Stat(src)
	if err != nil {
		return err
	}
	l.Debug.Log("Done.")

	l.Debug.Log("Changing file permissions at file '%v'.", des)
	if err := os.Chmod(des, f1.Mode()); err != nil {
		return err
	}
	l.Debug.Log("Done.")

	return nil
}

// pathExists checks to see if the given path exists, returns a bool
func pathExists(fp string) bool {
	l.Debug.Log("Does '%v' exist?", fp)
	_, err := os.Stat(fp)
	if err != nil {
		l.Debug.Log("No!")
		return false
	}
	l.Debug.Log("Yes!")
	return true
}

// splitPath splits the path string into the individual files and returns them in an array
func splitPath(path string, dirs bool) []string {
	l.Debug.Log("Splitting path '%v'.", path)
	sp := strings.Split(path, "/")
	if dirs {
		l.Debug.Log("Is a directory, so dropping off file name.")
		sp = sp[:len(sp)-1]
	}
	l.Debug.Log("Done.")

	return sp
}
