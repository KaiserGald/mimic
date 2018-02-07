// Package main
// 6 February 2018
// Code is licensed under the MIT License
// © 2018 Scott Isenberg

package main

import (
	"os"
	"testing"

	"github.com/KaiserGald/logger"
)

func TestMain(m *testing.M) {
	l = logger.New()
	processFlags()
	os.Exit(m.Run())
}

func TestHandleFlags(t *testing.T) {
	expSrc := "testsrc"
	expDes := "testdes"
	resSrc, resDes := handleFlags()
	if (resSrc != expSrc) && (resDes != expDes) {
		t.Errorf("Strings do not match. Expected ['%v', '%v'] got ['%v', '%v']", expSrc, expDes, resSrc, resDes)
	}
}
