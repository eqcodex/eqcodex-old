package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-yaml/yaml"
)

type Tracker struct {
	NpcId  int
	ZoneId int
	ItemId int
}

func loadTracker() (tracker *Tracker, err error) {
	inFile, err := ioutil.ReadFile("tracker.yml")
	if err != nil {
		err = fmt.Errorf("Failed to read tracker.yml: %s", err.Error())
		return
	}
	tracker = &Tracker{}
	if err = yaml.Unmarshal(inFile, &tracker); err != nil {
		err = fmt.Errorf("Failed to unmarshal: %s", err.Error())
		return
	}
	return
}

func saveTracker(tracker *Tracker) {
	var err error
	bData := []byte{}
	if bData, err = yaml.Marshal(&tracker); err != nil {
		log.Fatal("Failed to save tracker.yml:", err.Error())
	}
	ioutil.WriteFile("tracker.yml", bData, 0644)
}
