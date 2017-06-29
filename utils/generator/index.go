package main

import (
	//"html/template"
	"log"
	"os"
)

func generateIndex(instance *Instance) {
	var err error

	type ContentData struct {
		Name string
	}
	type PageData struct {
		Core    *CoreData
		Content *ContentData
	}

	page := &PageData{
		Core: getCore(),
		Content: &ContentData{
			Name: "Rawr",
		},
	}
	page.Core.Site.Title = "EQ Codex"

	t := getCoreTemplate(instance)

	t, err = t.ParseFiles(instance.yamlConfig.Templates + "index.tpl")
	if err != nil {
		log.Println("Failed to parse index: ", err.Error())
	}

	f, err := os.Create(instance.yamlConfig.Output + "index.html")
	if err != nil {
		log.Println("Failed to create file: ", err.Error())
	}

	err = t.Execute(f, page)
	if err != nil {
		log.Println("Failed to execute index: ", err.Error())
	}
	f.Close()
	return
}
