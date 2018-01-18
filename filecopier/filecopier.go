// Package filecopier
// 18 January 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package filecopier

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// CopyFile will copy the supplied file to the supplied destination
func CopyFile(srcfp string, desfp string) error {

	info, err := os.Stat(srcfp)
	if err != nil {
		return err
	}
	fmt.Println(info.IsDir())

	if p := pathExists(desfp); !p {
		fmt.Println(p)
		dirs := strings.Split(desfp, "/")
		dirs = dirs[:len(dirs)-1]
		var path string
		for i, dir := range dirs {
			if i == 0 {
				path = strings.Join([]string{path, dir}, "")
			} else {
				path = strings.Join([]string{path, dir}, "/")
			}
			if err = os.Mkdir(path, 0700); err != nil {
				return err
			}
		}
		fmt.Println(path)

	}

	from, err := os.Open(srcfp)
	if err != nil {
		log.Printf("Error opening source file: %v\n", err)
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(desfp, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Error opening destination file: %v\n", err)
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Printf("Error copying file: %v\n", err)
		return err
	}
	fmt.Println(from.Name())
	fmt.Println(to.Name())
	return nil
}

func pathExists(fp string) bool {
	_, err := os.Stat(fp)
	if err != nil {
		return false
	}
	return true
}
