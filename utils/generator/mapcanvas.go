package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/xackery/goeq/spawn"
)

type Line struct {
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
	X2 float64 `json:"x2"`
	Y2 float64 `json:"y2"`
}

type SpawnPoint struct {
	SpawnGroupID int     `json:"spawn_group_id"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Z            float64 `json:"z"`
}

/*
func generateMaps(instance *Instance) (err error) {
	zones := []*zone.Zone{}
	query := "SELECT * FROM zone WHERE min_status < 90 ORDER BY short_name"
	err = instance.db.Select(&zones, query)
	if err != nil {
		log.Println("Failed to select zones", err.Error())
		return
	}
	for _, zoneEntry := range zones {
		if zoneEntry.Short_name.String == "apprentice" ||
			zoneEntry.Short_name.String == "arena2" ||
			zoneEntry.Short_name.String == "arttest" ||
			zoneEntry.Short_name.String == "aviak" {

			//	zoneEntry.Short_name.String == "arttest" {
			continue
		}

		if err = mapByShortname(zoneEntry.Short_name.String, instance); err != nil {
			err = fmt.Errorf("Failed to generate map: %s", err.Error())
			return
		}
	}

	return
}*/

func mapCanvasByNpcId(npcid int64, shortname string, instance *Instance) (mapData string, err error) {
	type Index struct {
		Lines       []Line       `json:"lines"`
		SpawnPoints []SpawnPoint `json:"spawn_points"`
	}
	index := &Index{}

	lines, err := loadMap(shortname)
	if err != nil {
		fmt.Errorf("Failed to load map: %s", err.Error())
		return
	}
	index.Lines = lines

	rows, err := instance.db.Queryx(
		`SELECT spawn2.* FROM spawn2
		INNER JOIN spawnentry ON spawnentry.spawngroupid = spawn2.spawngroupid
		 WHERE spawnentry.npcid = ?
	 	 GROUP BY spawn2.id;`, npcid)
	if err != nil {
		fmt.Errorf("Error querying: %s", err.Error())
		return
	}

	for rows.Next() {
		spawnEntry := &spawn.Spawn2{}
		if err = rows.StructScan(&spawnEntry); err != nil {
			return
		}
		spawnData := SpawnPoint{
			X:            -spawnEntry.X,
			Y:            -spawnEntry.Y,
			Z:            spawnEntry.Z,
			SpawnGroupID: spawnEntry.Spawngroupid,
		}
		index.SpawnPoints = append(index.SpawnPoints, spawnData)
	}

	jsonData, err := json.Marshal(index)
	if err != nil {
		err = fmt.Errorf("Error marshallingL %s", err.Error())
	}
	mapData = string(jsonData)
	return
}

/*
func mapByShortname(shortname string, instance *Instance) (err error) {

	type Index struct {
		Lines       []Line       `json:"lines"`
		SpawnPoints []SpawnPoint `json:"spawn_points"`
	}
	index := &Index{}

	lines, err := loadMap(shortname)
	if err != nil {
		fmt.Errorf("Failed to load map: %s", err.Error())
		return
	}
	index.Lines = lines

	rows, err := instance.db.Queryx(
		`SELECT spawn2.* FROM spawnentry
		 INNER JOIN spawn2 ON spawnentry.spawngroupid = spawn2.spawngroupid
		 WHERE spawn2.zone = ?
	 	 GROUP BY spawn2.id;`, shortname)
	if err != nil {
		fmt.Errorf("Error querying: %s", err.Error())
		return
	}

	for rows.Next() {
		spawnEntry := &spawn.Spawn2{}
		if err = rows.StructScan(&spawnEntry); err != nil {
			return
		}
		spawnData := SpawnPoint{
			X:            -spawnEntry.X,
			Y:            -spawnEntry.Y,
			Z:            spawnEntry.Z,
			SpawnGroupID: spawnEntry.Spawngroupid,
		}
		index.SpawnPoints = append(index.SpawnPoints, spawnData)
	}

	err = json.NewEncoder(f).Encode(index)
	if err != nil {
		log.Println("Error requesting RestIndex:", err.Error())
	}
	return
}
*/
func loadMap(shortname string) (lines []Line, err error) {
	path := fmt.Sprintf("map/%s_1.txt", shortname)
	f, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("Error opening %s: %s", path, err.Error())
		return
	}

	bMap, err := ioutil.ReadAll(f)
	if err != nil {
		err = fmt.Errorf("Error finding map (%s): %s", shortname, err.Error())
		return
	}

	entries := strings.Split(string(bMap), "\n")

	for _, entry := range entries {
		record := strings.Split(entry, ",")
		if len(record) < 5 {
			continue
		}
		entries := strings.Split(record[0], " ")
		drawType := entries[0]
		if drawType == "L" {
			line := Line{}
			line.X1, _ = strconv.ParseFloat(strings.TrimSpace(entries[1]), 64)
			line.Y1, _ = strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
			line.X2, _ = strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
			line.Y2, _ = strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
			lines = append(lines, line)
		}

	}
	return
}
