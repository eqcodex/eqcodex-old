//Parses quests of zonedrops
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/goeq/npc"
	"github.com/xackery/goeq/zone"
)

var (
	config        *eqemuconfig.Config
	prefixPath    string
	db            Database
	lastZone      string
	questNoSpawns []string
	insertCount   int
)

type Database struct {
	instance *sqlx.DB
}

func main() {
	var err error
	if err = loadConfig(); err != nil {
		log.Fatal(err.Error())
	}

	db = Database{}
	if err = loadDatabase(&db); err != nil {
		log.Fatal(err.Error())
	}

	query := "DELETE FROM zone_drops WHERE is_quest_item = 1 OR is_quest_reward = 1"
	if _, err := db.instance.Exec(query); err != nil {
		log.Fatal(err.Error())
	}

	//os.Chdir("../../deploy/server/quests")
	prefixPath = "../../../deploy/server/quests/"
	prefixPath = "./quests/"
	err = filepath.Walk(prefixPath, visit)
	if err != nil {
		log.Fatal("Error filepath", err.Error())
	}
	fmt.Println("===============")
	for _, entry := range questNoSpawns {
		fmt.Println(entry)
	}
	fmt.Println("===============")
	fmt.Println("Created", insertCount, "entries")

}

func visit(path string, f os.FileInfo, err error) error {

	if f.IsDir() {
		//fmt.Println("Dir", path)
		return nil
	}

	dir, filename := filepath.Split(path)
	dir = strings.Replace(dir, prefixPath, "", -1)
	if len(dir) < 2 {
		return nil
	}
	dir = dir[0 : len(dir)-1]
	if !strings.Contains(filename, ".pl") && !strings.Contains(filename, ".lua") {
		return nil
	}

	if len(filename) < 1 {
		return nil
	}
	if lastZone != dir {
		lastZone = dir
		fmt.Println("Zone:", lastZone)
	}

	npcname := ""
	if strings.Contains(filename, ".pl") {
		npcname = filename[0 : len(filename)-3]
	}
	if strings.Contains(filename, ".lua") {
		npcname = filename[0 : len(filename)-4]
	}

	fmt.Printf("%s, ", filename)
	zone_id := getZoneByShortname(dir)
	if zone_id < 1 {
		//fmt.Println("Zone skipped", dir)
		return nil
	}

	npc_id := getNpcByNameOrId(npcname)
	if npc_id < 1 {
		//fmt.Println("NPC skipped", filename)
		questNoSpawns = append(questNoSpawns, fmt.Sprintf("%s/%s", dir, filename))
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parse := strings.TrimSpace(strings.Replace(strings.ToLower(scanner.Text()), " ", "", -1))

		isQuestReward := 0
		isQuestItem := 0

		var items []int
		if items = getQuestReward(parse); len(items) == 0 {
			if items = getQuestItem(parse); len(items) == 0 {
				continue
			}
			isQuestItem = 1
		} else {
			isQuestReward = 1
		}

		for _, item_id := range items {
			fmt.Println("Inserting", zone_id, npc_id, item_id, isQuestReward, isQuestItem, "for", dir, filename)
			insertQuery := "REPLACE INTO zone_drops (item_id, npc_id, zone_id, is_quest_reward, is_quest_item) VALUES (?, ?, ?, ?, ?)"
			stmt, _ := db.instance.Prepare(insertQuery)
			if _, err = stmt.Exec(item_id, npc_id, zone_id, isQuestReward, isQuestItem); err != nil {
				log.Fatal(err.Error())
			}
			insertCount++
			if insertCount%1000 == 0 {
				fmt.Println("Inserted", insertCount, "so far")
			}
		}

		//fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	//log.Fatal("Done")
	return nil
}

func getQuestReward(parse string) (items []int) {
	if !strings.Contains(parse, "summonitem(") {
		return
	}
	item_id := parse[strings.Index(parse, "summonitem(")+11:]
	if len(item_id) < 2 {
		return
	}
	item_id = item_id[0:strings.Index(item_id, ")")]
	if len(item_id) < 2 {
		return
	}
	if item_id == "0" {
		return
	}
	var err error
	var result int
	if result, err = strconv.Atoi(item_id); err != nil {
		return
	}
	items = append(items, result)
	return
}

func getQuestItem(parse string) (items []int) {
	var err error
	var result int
	if strings.Contains(parse, `check_handin(`) {

		//fmt.Println("1", parse)
		item_id := parse[strings.Index(parse, `check_handin(`)+13:]
		if len(item_id) < 2 {
			return
		}
		itemReg := regexp.MustCompile("([0-9]+)=>[0-9]+")
		strIds := itemReg.FindAllString(parse, -1)

		//fmt.Println("2", len(strIds))
		for _, strId := range strIds {
			//fmt.Println("3", strId)
			if strings.Index(strId, "=>") > 0 {
				strId = strings.Split(strId, "=>")[0]
			}

			//fmt.Println("4", strId)
			if result, err = strconv.Atoi(strId); err != nil {
				continue
			}
			items = append(items, result)
		}
		return
	} else if strings.Contains(parse, "item_lib.check_turn_in") {
		itemReg := regexp.MustCompile("item[0-9]=([0-9]+)")
		//fmt.Println(itemReg.FindAll(b, n) (parse))
		//fmt.Println("1", parse)
		strIds := itemReg.FindAllString(parse, -1)
		//fmt.Println("2", strIds)
		for _, strId := range strIds {
			//fmt.Println("3", strId)
			if strings.Index(strId, "=") > 0 {
				strId = strings.Split(strId, "=")[1]
			}
			//fmt.Println("4", strId)
			if result, err = strconv.Atoi(strId); err != nil {
				continue
			}
			items = append(items, result)
		}
	}
	return
}
func loadConfig() error {
	if config != nil {
		return nil
	}
	var err error
	if config, err = eqemuconfig.GetConfig(); err != nil {
		return err
	}
	return nil
}

func loadDatabase(db *Database) error {
	if db == nil {
		return fmt.Errorf("No database")
	}
	var err error
	if db.instance, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Db)); err != nil {
		return fmt.Errorf("error connecting to db: %s", err.Error())
	}
	return nil
}

func getZoneByShortname(shortname string) int {
	zones := []zone.Zone{}
	query := "SELECT * from zone WHERE min_status < 90 AND short_name = ? LIMIT 1" //was min_status = 0
	if err := db.instance.Select(&zones, query, shortname); err != nil {
		log.Fatal("Error getting ", shortname, ": ", err.Error())
		return 0
	}
	if len(zones) == 0 {
		return 0
	}
	return zones[0].Zoneidnumber
}

func getNpcByNameOrId(npcname string) int64 {
	npcs := []npc.NpcTypes{}
	query := `SELECT npc.* 
	FROM npc_types npc 
	INNER JOIN spawnentry ON spawnentry.npcid = npc.id
	INNER JOIN spawn2 ON spawn2.spawngroupid = spawnentry.spawngroupid
	WHERE npc.id = ? OR npc.Name = ?`
	if err := db.instance.Select(&npcs, query, npcname, npcname); err != nil {
		log.Fatal("Error getting", npcname, err.Error())
		return 0
	}
	if len(npcs) == 0 {
		return 0
	}
	return npcs[0].Id.Int64
}
