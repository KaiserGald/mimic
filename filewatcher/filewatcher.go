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

	"github.com/KaiserGald/filewatcher/filecopier"
	"github.com/radovskyb/watcher"
)

// WatchFiles will watch the files at the specified filepath and will fire off an
// event when a change happens
func WatchFiles(srcfp string, desfp string) error {
	w := watcher.New()
	w.IgnoreHiddenFiles(true)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	relfp := strings.Join([]string{dir, srcfp}, "/")
	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Println(event) // Print the event's info.
				rel := strings.Replace(event.Path, relfp, "", -1)
				des := strings.Join([]string{desfp, rel}, "")
				src := strings.Join([]string{srcfp, rel}, "")
				fmt.Println(src, des)
				err := filecopier.CopyFile(src, des)
				if err != nil {
					log.Printf("Error copying file: %v\n", err)
					return
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
