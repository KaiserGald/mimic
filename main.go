// package main
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/KaiserGald/filewatcher/filewatcher"
)

var (
	watch string
)

func processFlags() (string, string) {
	flag.StringVar(&watch, "w", "", "watches the specified files and copies them to the specified location. Example: fw -w SOURCE:DESTINATION")
	flag.StringVar(&watch, "watch", "", "watches the specified files and copies them to the specified location. Example: fw -w SOURCE:DESTINATION")

	flag.Parse()

	src, des := handleFlags()
	return src, des
}

func handleFlags() (string, string) {
	fps := strings.Split(watch, ":")
	src := fps[0]
	des := fps[1]
	return src, des
}

func main() {
	srcfp, desfp := processFlags()
	fmt.Println(srcfp, desfp)
	err := filewatcher.WatchFiles(srcfp, desfp)
	if err != nil {
		log.Fatalln("Error starting filewatcher: %v\n", err)
	}

	waitForSignal()
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch

	log.Printf("Got signal: %v, exiting.\n", s)
	time.Sleep(2 * time.Second)
}
