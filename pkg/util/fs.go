package util

import "os"

func Exists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
