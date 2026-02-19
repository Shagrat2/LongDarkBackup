package main

import (
	"os"
	"strings"
)

func FileExists(Dir string, FilaPref string) bool {

	items, err := os.ReadDir(Dir)
	if err != nil {
		return false
	}

	for _, itm := range items {
		if strings.HasPrefix(itm.Name(), FilaPref) {
			return true
		}
	}

	return false
}
