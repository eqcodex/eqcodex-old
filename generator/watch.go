package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"time"
)

func watchForChanges(instance *Instance) {
	var err error
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	generateTemplates(instance)
	templateTrigger := make(chan bool)
	go func() {
		nextRun := time.Now().Add(1 * time.Second)
		for {
			<-templateTrigger
			if time.Now().Before(nextRun) {
				continue
			}
			generateTemplates(instance)
			nextRun = time.Now().Add(1 * time.Second)
		}
	}()
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
				templateTrigger <- true
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	if err = watcher.Watch(instance.yamlConfig.Templates); err != nil {
		log.Fatal(err)
	}
	if err = watcher.Watch(instance.yamlConfig.Templates + "item/"); err != nil {
		log.Fatal(err)
	}
}
