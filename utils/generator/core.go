package main

import (
	"fmt"
	"html/template"
	"log"
)

type CoreData struct {
	Site *SiteData
}

type SiteData struct {
	Title   string
	Name    string
	FavIcon string
}

var core *CoreData

func getCore() *CoreData {
	if core != nil {
		return core
	}
	core = &CoreData{
		Site: &SiteData{
			Title:   "EQCodex",
			Name:    "Everquest Codex",
			FavIcon: "favicon.ico",
		},
	}
	return core
}

func getCoreTemplate(instance *Instance) (t *template.Template) {
	t, err := template.ParseFiles(instance.yamlConfig.Templates + "core.tpl")
	if err != nil {
		log.Println("Failed to parse core: ", err.Error())
	}
	return t
}

func showPercent(message string, cur int, max int, remaining string, color string) {
	switch color {
	case "black":
		color = "\033[30m"
	case "maroon":
		color = "\033[31m"
	case "green":
		color = "\033[32m"
	case "yellow":
		color = "\033[33m"
	case "blue":
		color = "\033[34m"
	case "purple":
		color = "\033[35m"
	case "cyan":
		color = "\033[36m"
	case "white":
		color = "\033[37m"
	case "gray":
		color = "\033[38m"
	case "red":
		color = "\033[39m"
	default:
		color = "\033[0m"
	}

	dotCount := 30
	fmt.Printf("\033[2K\033[120D") //first of line
	val := float64(cur) / float64(max) * float64(dotCount)
	fmt.Printf("%s%s - [", color, message)
	for i := 0; i < dotCount; i++ {
		if int(val) >= i {
			fmt.Printf(".")
		} else {
			fmt.Printf(" ")
		}
	}
	fmt.Printf("] - %s\033[0m", remaining)
	if cur == max {
		fmt.Printf("\n")
	}
}
