// package main
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/KaiserGald/logger"
	"github.com/KaiserGald/mimic/filewatcher"
	"github.com/logrusorgru/aurora"
)

var (
	watch   string
	color   bool
	dev     bool
	verbose bool
	quiet   bool
	l       *logger.Logger
	au      aurora.Aurora
)

func processFlags() (string, string) {
	flag.BoolVar(&color, "c", false, "Short version of -color. Starts mimic with colored output.")
	flag.BoolVar(&color, "color", false, "Starts mimic with colored output.")

	flag.BoolVar(&dev, "d", false, "Short version of -dev. Starts mimic in dev mode.")
	flag.BoolVar(&dev, "dev", false, "Starts mimic in dev mode.")

	flag.BoolVar(&quiet, "q", false, "Short version of -quiet. Starts mimic with quiet output.")
	flag.BoolVar(&quiet, "quiet", false, "Starts mimic with quiet output.")

	flag.BoolVar(&verbose, "v", false, "Short version of -verbose. Starts mimic with verbose output.")
	flag.BoolVar(&verbose, "verbose", false, "Starts mimic with verbose output.")

	flag.StringVar(&watch, "w", "", "Short version of -watch. Watches the specified files and copies them to the specified location. Example: mimic -w 'SOURCE:DESTINATION'")
	flag.StringVar(&watch, "watch", "", "Watches the specified files and copies them to the specified location. Example: mimic -watch 'SOURCE:DESTINATION'")

	flag.Parse()

	src, des := handleFlags()
	return src, des
}

func handleFlags() (string, string) {
	l.ShowColor(color)
	au = aurora.NewAurora(color)
	var src, des string
	if watch != "" {
		fps := strings.Split(watch, ":")
		src = fps[0]
		des = fps[1]
	} else {
		fmt.Printf("\n%v needs to have a source and destination directory supplied via the %v or %v flag. Usage is: %v %v %v%v%v%v%v.\nThe source directory must already exist. %v will automatically create the destination directories and clone any existing files\nfrom the %v directory into the %v directory.\n\n", au.Magenta("Mimic"), au.Cyan("-w"), au.Cyan("-watch"), au.Gray("mimic"), au.Cyan("-w"), au.Gray("'"), au.Red("SOURCE"), au.Gray(":"), au.Green("DESTINATION"), au.Gray("'"), au.Magenta("Mimic"), au.Red("source"), au.Green("destination"))
		usage()
		l.Notice.Log("Exiting now...")
		os.Exit(0)
	}

	if quiet {
		l.SetLogLevel(logger.ErrorsOnly)
	}
	if verbose {
		l.SetLogLevel(logger.Verbose)
	}
	if dev {
		l.SetLogLevel(logger.All)
	}
	return src, des
}

func main() {
	l = logger.New()
	srcfp, desfp := processFlags()
	l.Info.Log("Starting filewatcher...")
	err := filewatcher.WatchFiles(srcfp, desfp, l)
	if err != nil {
		l.Error.Log("Error starting filewatcher: %v", err)
	}

}

func usage() {
	fmt.Printf("%v %v%v\n", au.Gray("Usage of"), au.Magenta("mimic"), au.Gray(":"))
	fmt.Printf("\t%v,%v\n\t\tStarts mimic with colored output.\n", au.Cyan("-c"), au.Cyan("-color"))
	fmt.Printf("\t%v,%v\n\t\tStarts mimic in dev mode.\n", au.Cyan("-d"), au.Cyan("-dev"))
	fmt.Printf("\t%v,%v\n\t\tStarts mimic in quiet output mode.\n", au.Cyan("-q"), au.Cyan("-quiet"))
	fmt.Printf("\t%v,%v\n\t\tStarts mimic in verbose output mode.\n", au.Cyan("-v"), au.Cyan("-verbose"))
	fmt.Printf("\t%v,%v string\n\t\tWatches the specified files and copies them to the specified location. Example: %v %v %v%v%v%v%v\n", au.Cyan("-w"), au.Cyan("-watch"), au.Gray("mimic"), au.Cyan("-w"), au.Gray("'"), au.Red("SOURCE"), au.Gray(":"), au.Green("DESTINATION"), au.Gray("'"))
}
