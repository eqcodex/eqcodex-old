package main

import (
	"regexp"
	"strings"
)

func cleanName(name string) (cleanName string) {
	var re = regexp.MustCompile(`[^0-9A-Za-z_]+`)
	cleanName = strings.Replace(name, " ", "_", -1)
	cleanName = strings.Replace(cleanName, "#", "", -1)
	cleanName = strings.TrimSpace(re.ReplaceAllString(cleanName, ""))
	cleanName = strings.Replace(cleanName, "_", " ", -1)
	return
}

func cleanUrl(srcUrl string) (cleanUrl string) {
	cleanUrl = strings.Replace(cleanName(srcUrl), " ", "-", -1)
	return
}

func getCategory(slots int) string {
	if (slots&8192) == 8192 || //primary
		(slots&16384) == 16384 || //secondary
		(slots&2048) == 2048 { //range
		return "Weapon"
	}
	if (slots&1) == 1 ||
		(slots&4) == 4 ||
		(slots&8) == 8 ||
		(slots&16) == 16 ||
		(slots&32) == 32 ||
		(slots&64) == 64 ||
		(slots&128) == 128 ||
		(slots&256) == 256 ||
		(slots&1536) == 1536 ||
		(slots&4096) == 4096 ||
		(slots&98304) == 98304 ||
		(slots&131072) == 131072 ||
		(slots&262144) == 262144 ||
		(slots&524288) == 524288 ||
		(slots&1048576) == 1048576 {
		return "Gear"
	}
	return "Item"
}
