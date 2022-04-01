package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
)

func defaultConfigPath() string {
	// get binary path
	bp := getBinaryPath()
	if dirIsWriteable(bp) == nil {
		return bp
	}

	// binary path is not write-able, use alternative path
	cp := configdir.LocalConfig("nabarr")
	if _, err := os.Stat(cp); os.IsNotExist(err) {
		if e := os.MkdirAll(cp, os.ModePerm); e != nil {
			panic("failed to create nabarr config directory")
		}
	}

	return cp
}

func getBinaryPath() string {
	// get current binary path
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		// get current working dir
		if dir, err = os.Getwd(); err != nil {
			panic("failed to determine current binary location")
		}
	}

	return dir
}

func dirIsWriteable(dir string) error {
	// credits: https://stackoverflow.com/questions/20026320/how-to-tell-if-folder-exists-and-is-writable
	info, err := os.Stat(dir)
	if err != nil {
		return errors.New("path does not exist")
	}

	if !info.IsDir() {
		return errors.New("path is not a directory")
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		fmt.Println("Write permission bit is not set on this file for user")
		return errors.New("write permission not set for user")
	}

	return nil
}
