// package main
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/radovskyb/watcher"
)

func main() {
	w := watcher.New()
	w.IgnoreHiddenFiles(true)

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive("."); err != nil {
		log.Fatalln(err)
	}

	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	fmt.Println()

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
