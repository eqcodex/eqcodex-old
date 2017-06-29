package main

import (
	"strings"
)

func cleanName(name string) (cleanName string) {
	cleanName = strings.Replace(name, "_", " ", -1)
	cleanName = strings.Replace(cleanName, "#", "", -1)
	return
}
