package goutils

import "os"

func PathExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
