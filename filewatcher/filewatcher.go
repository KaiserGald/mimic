// Package filewatcher
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package filewatcher

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/KaiserGald/logger"
	"github.com/KaiserGald/mimic/filehandler"
	"github.com/radovskyb/watcher"
)

var l *logger.Logger

// initWatcher will initialize the watcher with any configuration an return the watcher, it also gets and returns the relative filepath to the source directory
func initWatcher(srcfp string) (*watcher.Watcher, string, error) {
	w := watcher.New()
	w.IgnoreHiddenFiles(true)
	filehandler.Init(l)

	// get relative file path
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, "", err
	}
	relfp := strings.Join([]string{dir, srcfp}, "/")

	return w, relfp, nil
}

// WatchFiles will watch the files at the specified filepath and will fire off an
// event when a change happens
func WatchFiles(srcfp, desfp string, lg *logger.Logger) error {
	l = lg
	l.Debug.Log("Initializing watcher...")
	w, relfp, err := initWatcher(srcfp)
	if err != nil {
		return err
	}
	l.Debug.Log("Done")
	l.Notice.Log("Initializing the destination file tree...")
	err = initializeFileTree(srcfp, desfp, relfp)
	if err != nil {
		return err
	}
	l.Debug.Log("Done initializing destination file tree.")
	// listen for events
	l.Info.Log("Listening for events at '%v'.", relfp)
	go func() {
		for {
			select {
			case event := <-w.Event:
				l.Debug.Log(event.String())
				switch event.Op.String() {
				case "CREATE":
					l.Debug.Log("CREATE event occured at '%v'", event.Path)
					err = handleCreate(event, srcfp, desfp, relfp)
					if err != nil {
						return
					}
					l.Debug.Log("CREATE event handled.")

				case "WRITE":
					l.Debug.Log("WRITE event occured at '%v'", event.Path)
					err = handleWrite(event, srcfp, desfp, relfp)
					if err != nil {
						return
					}
					l.Debug.Log("WRITE event handled.")

				case "REMOVE":
					l.Debug.Log("REMOVE event occured at '%v'", event.Path)
					err = handleRemove(event, srcfp, desfp, relfp)
					if err != nil {
						return
					}
					l.Debug.Log("REMOVE event handled.")

				case "RENAME":
					l.Debug.Log("RENAME event occured at '%v'", event.Path)
					err = handleRename(event, desfp, relfp)
					if err != nil {
						return
					}
					l.Debug.Log("RENAME event handled.")

				case "CHMOD":
					l.Debug.Log("CHMOD event occured at '%v'", event.Path)
					err = handleChmod(event, srcfp, desfp, relfp)
					if err != nil {
						return
					}
					l.Debug.Log("CHMOD event handled.")

				case "MOVE":
					l.Debug.Log("MOVE event occured at '%v'", event.Path)
					err = handleMove(event, srcfp, desfp, relfp)
					if err != nil {
						return
					}
					l.Debug.Log("MOVE event handled.")
				}

			case err := <-w.Error:
				l.Error.Log(err.Error())
				return
			case <-w.Closed:
				return
			}
		}
	}()

	l.Debug.Log("Adding '%v' to be watched...", srcfp)
	if err := w.AddRecursive(srcfp); err != nil {
		return err
	}

	for path, f := range w.WatchedFiles() {
		l.Info.Log("%s is mimicking %s", f.Name(), path)
	}

	l.Debug.Log("Done.")
	l.Notice.Log("Mimic successfully started!")

	if err := w.Start(time.Millisecond * 100); err != nil {
		return err
	}

	return nil
}

func initializeFileTree(srcfp, desfp, relfp string) error {
	l.Debug.Log("Mapping source tree in '%v'...", srcfp)
	tree, err := mapTree(srcfp)
	if err != nil {
		return err
	}
	l.Debug.Log("Tree: %v", tree)
	l.Debug.Log("Done.")
	l.Debug.Log("Starting to copy source tree to destination tree...")

	for file := range tree {

		l.Debug.Log("file: %v", file)

		l.Debug.Log("Building file paths...")
		src, des := buildPaths("/"+file, srcfp, desfp, relfp)
		l.Debug.Log("Done.")

		l.Debug.Log("Is file a directory?")
		if !tree[file].IsDir() {
			l.Debug.Log("No!")
			l.Info.Log("Copying '%v' into '%v'", src, des)
			err := filehandler.CopyFile(src, des)
			if err != nil {
				l.Error.Log("Error copying a file: %v", err)
				return err
			}
		} else {
			l.Debug.Log("Yes!")
			l.Info.Log("Copying '%v' into '%v'", src, des)
			err := filehandler.CopyDir(src, des)
			if err != nil {
				l.Error.Log("Error copying a directory: %v", err)
				return err
			}
		}
		l.Debug.Log("Copy complete!")
	}

	l.Debug.Log("Done copying source tree to destination tree.")
	return nil
}

// handleCreate handles the create events for both directories and files.
func handleCreate(event watcher.Event, srcfp, desfp, relfp string) error {
	if event.IsDir() {
		l.Debug.Log("Building paths...")
		src, des := buildPaths(event.Path, srcfp, desfp, relfp)
		l.Debug.Log("Done.")
		l.Info.Log("Copying directory %v to %v...", src, des)
		err := filehandler.CopyDir(src, des)
		if err != nil {
			l.Error.Log("Error copying directory: %v", err)
			return err
		}
		l.Debug.Log("Done copying directory.")
	} else {
		l.Debug.Log("Building paths...")
		src, des := buildPaths(event.Path, srcfp, desfp, relfp)
		l.Debug.Log("Done.")
		l.Info.Log("Copying file %v to %v...", src, des)
		err := filehandler.CopyFile(src, des)
		if err != nil {
			l.Error.Log("Error copying file: %v", err)
			return err
		}
		l.Debug.Log("Done copying file.")
	}
	return nil
}

// handleWrite handles the write events for files.
func handleWrite(event watcher.Event, srcfp, desfp, relfp string) error {
	l.Debug.Log("Is WRITE event at a direcory?")
	if !event.IsDir() {
		l.Debug.Log("No!")
		l.Debug.Log("Building paths...")
		src, des := buildPaths(event.Path, srcfp, desfp, relfp)
		l.Debug.Log("Done.")
		l.Info.Log("Copying '%v' into '%v'.", src, des)
		err := filehandler.CopyFile(src, des)
		if err != nil {
			l.Error.Log("Error copying file: %v", err)
			return err
		}
		l.Debug.Log("Done copying file.")
	} else {
		l.Debug.Log("Yes...")
	}
	return nil
}

// handleRemove handles the remove events for files
func handleRemove(event watcher.Event, srcfp, desfp, relfp string) error {
	l.Debug.Log("Building path...")
	_, des := buildPaths(event.Path, srcfp, desfp, relfp)
	l.Debug.Log("Done.")
	l.Info.Log("Removing '%v'.", des)
	err := filehandler.Remove(des)
	if err != nil {
		l.Error.Log("Error deleting file: %v", err)
		return err
	}
	l.Debug.Log("Done.")

	return nil
}

func handleRename(event watcher.Event, desfp, relfp string) error {
	l.Debug.Log("Building paths...")
	path := strings.Split(event.Path, " -> ")
	l.Debug.Log("path: %v", path)
	rel := strings.Replace(path[1], relfp, "", -1)
	l.Debug.Log("rel: %v", rel)
	old := strings.Join([]string{desfp, event.Name()}, "/")
	l.Debug.Log("old: %v", old)
	new := strings.Join([]string{desfp, rel}, "")
	l.Debug.Log("new: %v", new)
	l.Debug.Log("Done.")
	l.Info.Log("Renaming '%v' to '%v'.", old, new)
	err := filehandler.Rename(old, new)
	if err != nil {
		l.Error.Log("Error renaming file: %v", err)
		return err
	}
	l.Debug.Log("Done.")
	return nil
}

func handleChmod(event watcher.Event, srcfp, desfp, relfp string) error {
	l.Debug.Log("Building paths...")
	src, des := buildPaths(event.Path, srcfp, desfp, relfp)
	l.Debug.Log("Done.")
	l.Info.Log("Copying file permissions from '%v' to '%v'.", src, des)
	err := filehandler.Chmod(src, des)
	if err != nil {
		l.Error.Log("Error changing permissions: %v", err)
	}
	l.Debug.Log("Done.")
	return nil
}

func handleMove(event watcher.Event, srcfp, desfp, relfp string) error {
	l.Debug.Log("Building paths...")
	path := strings.Split(event.Path, " -> ")
	l.Debug.Log("path: %v", path)
	path[0] = strings.Replace(path[0], relfp, "", -1)
	path[1] = strings.Replace(path[1], relfp, "", -1)
	l.Debug.Log("Move Source Path: %v", path[0])
	l.Debug.Log("Move Destination Path: %v", path[1])
	src, _ := buildPaths(path[0], desfp, desfp, relfp)
	_, des := buildPaths(path[1], srcfp, desfp, relfp)
	l.Debug.Log("Is source directory?")
	if event.IsDir() {
		l.Debug.Log("Yes!")
		l.Info.Log("Moving '%v' to '%v'.", src, des)
		err := filehandler.CopyDir(src, des)
		if err != nil {
			l.Error.Log("Error moving directory: %v\n", err)
		}
		l.Debug.Log("Done moving directory.")
	} else {
		l.Debug.Log("No!")
		l.Info.Log("Moving '%v' to '%v'", src, des)
		err := filehandler.CopyFile(src, des)
		if err != nil {
			l.Error.Log("Error moving file: %v\n", err)
		}
		l.Debug.Log("Done moving file.")
	}
	l.Debug.Log("Removing source file...")
	err := filehandler.Remove(src)
	if err != nil {
		l.Error.Log("Error removing source file: %v\n", err)
	}
	l.Debug.Log("Done.")
	return nil
}

// copyFile takes the file that triggered the event and copies it to the destination
func buildPaths(ep string, srcfp string, desfp string, relfp string) (string, string) {
	rel := strings.Replace(ep, relfp, "", -1)
	l.Debug.Log("rel: %v", rel)
	des := strings.Join([]string{desfp, rel}, "")
	l.Debug.Log("des: %v", des)
	src := strings.Join([]string{srcfp, rel}, "")
	l.Debug.Log("src: %v", src)
	return src, des
}

// mapTree returns a map of the file tree being watched
func mapTree(name string) (map[string]os.FileInfo, error) {

	tree := make(map[string]os.FileInfo)
	l.Debug.Log("tree: %v", tree)

	return tree, filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files := strings.Split(path, "/")
		l.Debug.Log("files: %v", files)
		files = files[1:]
		l.Debug.Log("files: %v", files)
		path = strings.Join(files, "/")
		l.Debug.Log("path: %v", path)
		l.Debug.Log("len(path): %v", len(path))
		if len(path) != 0 {
			tree[path] = info
		}
		return nil
	})
}
