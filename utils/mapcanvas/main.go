package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/goeq/spawn"
	"github.com/xackery/goeq/zone"
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

type Instance struct {
	yamlConfig  *YamlConfig
	eqemuConfig *eqemuconfig.Config
	db          *sqlx.DB
}

func main() {
	var err error
	instance := &Instance{}
	if err = initialize(instance); err != nil {
		log.Fatal(err)
	}

	if err = generateMaps(instance); err != nil {
		log.Fatal(err)
	}
}

func initialize(instance *Instance) (err error) {
	//Load Config
	instance.yamlConfig, err = loadYamlConfig()
	if err != nil {
		log.Fatal("Error while loading yaml: ", err.Error())
	}

	//Load EQEMU Config
	instance.eqemuConfig, err = loadEqemuConfig()
	if err != nil {
		log.Fatal("Error while loading yaml: ", err.Error())
	}

	//Connect to DB
	instance.db = &sqlx.DB{}
	if instance.db, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", instance.eqemuConfig.Database.Username, instance.eqemuConfig.Database.Password, instance.eqemuConfig.Database.Host, instance.eqemuConfig.Database.Port, instance.eqemuConfig.Database.Db)); err != nil {
		log.Fatalf("Database error: %s", err.Error())
	}
	return
}

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
}

func mapByShortname(shortname string, instance *Instance) (err error) {

	type Index struct {
		Lines       []Line       `json:"lines"`
		SpawnPoints []SpawnPoint `json:"spawn_points"`
	}
	index := &Index{}
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

	reader := csv.NewReader(strings.NewReader(string(bMap)))
	records, err := reader.ReadAll()
	if err != nil {
		err = fmt.Errorf("Error reading map (%s): %s", shortname, err.Error())
		return
	}

	for _, record := range records {
		entries := strings.Split(record[0], " ")
		drawType := entries[0]
		if drawType == "L" {
			line := Line{}
			line.X1, _ = strconv.ParseFloat(strings.TrimSpace(entries[1]), 64)
			line.Y1, _ = strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
			line.X2, _ = strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
			line.Y2, _ = strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
			index.Lines = append(index.Lines, line)
		}
	}

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

	os.MkdirAll(instance.yamlConfig.Output, 0744)
	f.Close()
	f, err = os.Create(instance.yamlConfig.Output + shortname + ".json")
	if err != nil {
		log.Println("Failed to create file:", err.Error())
		return
	}

	err = json.NewEncoder(f).Encode(index)
	if err != nil {
		log.Println("Error requesting RestIndex:", err.Error())
	}
	return
}
