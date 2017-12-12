package main

import (
	"log"

	"github.com/go-fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

func main() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("Change detected! %#v\n", event)

			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	if err := watcher.Add("test/file.txt"); err != nil {
		log.Fatal(err)
	}

	<-done
}
