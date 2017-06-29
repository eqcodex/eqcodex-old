package main

import (
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
