// Package filewatcher
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package filewatcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/KaiserGald/filewatcher/filehandler"
	"github.com/radovskyb/watcher"
)

// WatchFiles will watch the files at the specified filepath and will fire off an
// event when a change happens
func WatchFiles(srcfp string, desfp string) error {
	w := watcher.New()
	w.IgnoreHiddenFiles(true)

	// get relative file path
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	relfp := strings.Join([]string{dir, srcfp}, "/")

	// listen for events
	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Println(event) // Print the event's info.
				fmt.Println(event.Name())

				switch event.Op.String() {
				case "CREATE":
					fmt.Println("create happened")
					if event.IsDir() {
						src, des := buildPaths(event.Path, srcfp, desfp, relfp)
						err = filehandler.CopyDir(src, des)
						if err != nil {
							log.Printf("Error copying directory: %v\n", err)
						}
					} else {
						src, des := buildPaths(event.Path, srcfp, desfp, relfp)
						err = filehandler.CopyFile(src, des)
						if err != nil {
							log.Printf("Error copying file: %v\n", err)
						}
					}

				case "WRITE":
					fmt.Println("write happened")
					if !event.IsDir() {
						src, des := buildPaths(event.Path, srcfp, desfp, relfp)
						err = filehandler.CopyFile(src, des)
						if err != nil {
							log.Printf("Error copying file: %v\n", err)
						}
					}
				case "REMOVE":
					_, des := buildPaths(event.Path, srcfp, desfp, relfp)
					err = filehandler.Remove(des)
					if err != nil {
						fmt.Println("Error deleting file: %v\n", err)
					}

				case "RENAME":
					path := strings.Split(event.Path, " -> ")
					rel := strings.Replace(path[1], relfp, "", -1)
					old := strings.Join([]string{desfp, event.Name()}, "/")
					new := strings.Join([]string{desfp, rel}, "")
					err = filehandler.Rename(old, new)
					if err != nil {
						fmt.Println("Error renaming file: %v\n", err)
					}
				case "CHMOD":
					src, des := buildPaths(event.Path, srcfp, desfp, relfp)
					err = filehandler.Chmod(src, des)
					if err != nil {
						fmt.Println("Error changing permissions: %v\n", err)
					}

				case "MOVE":
					path := strings.Split(event.Path, " -> ")
					path[0] = strings.Replace(path[0], relfp, "", -1)
					path[1] = strings.Replace(path[1], relfp, "", -1)
					fmt.Printf("Path 0: %v\n", path[0])
					fmt.Printf("Path 1: %v\n", path[1])
					src, _ := buildPaths(path[0], desfp, desfp, relfp)
					_, des := buildPaths(path[1], srcfp, desfp, relfp)
					fmt.Println(src, des)
					if event.IsDir() {
						err = filehandler.CopyDir(src, des)
						if err != nil {
							fmt.Println("Error copying directory: %v\n", err)
						}
					} else {
						err = filehandler.CopyFile(src, des)
						if err != nil {
							fmt.Println("Error copying file: %v\n", err)
						}
					}
					err = filehandler.Remove(src)
					if err != nil {
						fmt.Println("Error removing source file: %v\n", err)
					}
				}

			case err := <-w.Error:
				log.Println(err)
				return
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(srcfp); err != nil {
		return err
	}

	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	fmt.Println()

	if err := w.Start(time.Millisecond * 100); err != nil {
		return err
	}

	return nil
}

// copyFile takes the file that triggered the event and copies it to the destination
func buildPaths(ep string, srcfp string, desfp string, relfp string) (string, string) {
	rel := strings.Replace(ep, relfp, "", -1)
	des := strings.Join([]string{desfp, rel}, "")
	src := strings.Join([]string{srcfp, rel}, "")
	return src, des
}

func initDir() {

}
