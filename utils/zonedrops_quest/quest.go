//Parses quests of zonedrops
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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
	recordCount   int
	dirCount      int
	insertTotal   int
	npc_id        int64
	zone_id       int
)

type Database struct {
	instance *sqlx.DB
}

type QuestData struct {
	tag         string
	title       string
	description string
	stepnum     int   `db:"quest_stepnum"`
	questid     int64 `db:"quest_id"`
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

	fmt.Printf("\033[48;5;7m") //white background
	fmt.Printf("\033[38;5;0m") //black text
	fmt.Printf("Scanning quests dir for items to put into zone_drops...")
	fmt.Printf("\033[48;5;0m\033[38;5;0m\033[0m\n") //restore text
	//os.Chdir("../../deploy/server/quests")
	//prefixPath = "../../../deploy/server/quests/"
	prefixPath = "quests/"
	err = filepath.Walk(prefixPath, visit)
	if err != nil {
		log.Fatal("Error filepath: ", err.Error())
	}
	fmt.Println("===============")
	for _, entry := range questNoSpawns {
		fmt.Println(entry)
	}
	fmt.Println("===============")
	fmt.Println("Created", insertTotal, " total entries")

}

func visit(path string, f os.FileInfo, err error) error {
	if f == nil {
		return fmt.Errorf("File does not exist")
	}
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
		recordCount = 0
		insertCount = 0
		files, _ := ioutil.ReadDir(prefixPath + dir)
		dirCount = len(files)
		//fmt.Println("Zone:", lastZone)
	}
	recordCount++
	showPercent(dir, recordCount, dirCount, insertCount, "green")
	//log.Fatal("Done", prefixPath, " ", dir, " ", recordCount, dirCount)
	npcname := ""
	if strings.Contains(filename, ".pl") {
		npcname = filename[0 : len(filename)-3]
	}
	if strings.Contains(filename, ".lua") {
		npcname = filename[0 : len(filename)-4]
	}

	//fmt.Printf("%s, ", filename)
	zone_id = getZoneByShortname(dir)
	if zone_id < 1 {
		//fmt.Println("Zone skipped", dir)
		return nil
	}

	npc_id = getNpcByNameOrId(npcname)
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
	isIgnore := false
	questData := &QuestData{}
	for scanner.Scan() {
		parse := strings.TrimSpace(strings.Replace(strings.ToLower(scanner.Text()), " ", "", -1))

		if !isIgnore {
			isIgnore = ignoreParse(parse)
		}

		if isIgnore {
			endIgnore := endParse(parse, questData)
			if !endIgnore {
				continue
			}
		}

		//line found to contain quest data
		if questParse(parse, questData) {
			continue
		}

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
			insertCount++
			insertTotal++
			//fmt.Println("Inserting", zone_id, npc_id, item_id, isQuestReward, isQuestItem, "for", dir, filename)
			insertQuery := "REPLACE INTO zone_drops (item_id, npc_id, zone_id, is_quest_reward, is_quest_item) VALUES (?, ?, ?, ?, ?)"
			stmt, _ := db.instance.Prepare(insertQuery)
			if _, err = stmt.Exec(item_id, npc_id, zone_id, isQuestReward, isQuestItem); err != nil {
				log.Fatal(err.Error())
			}
			//if insertCount%1000 == 0 {
			//	fmt.Println("Inserted", insertCount, "so far")
			//}
		}

		//fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	//log.Fatal("Done")
	return nil
}

func questParse(parse string, questData *QuestData) bool {
	var err error
	if strings.Contains(parse, "#!define") || strings.Contains(parse, "--!define") {
		qData := strings.Split(parse, "|")
		if len(qData) != 4 {
			log.Print("Failed to parse line", parse, "due to split of | being != 4 (", len(qData), ")\n")
			return true
		}

		questData.tag = qData[1]
		if len(questData.tag) > 64 {
			questData.tag = questData.tag[0:64]
		}
		questData.title = qData[2]
		if len(questData.title) > 128 {
			questData.title = questData.title[0:128]
		}
		questData.description = qData[3]
		if len(questData.description) > 256 {
			questData.description = questData.description[0:256]
		}

		insertQuery := "REPLACE INTO quests (tag, title, description) VALUES (?, ?, ?)"
		stmt, _ := db.instance.Prepare(insertQuery)
		if _, err = stmt.Exec(questData.tag, questData.title, questData.description); err != nil {
			log.Fatal(err.Error())
		}
		stmt.Close()
		if err = db.instance.Get(questData.questid, "SELECT id FROM quests WHERE tag = ?", questData.tag); err != nil {
			log.Fatal(err.Error())
		}

		return true
	}

	if strings.Contains(parse, "#!quest") || strings.Contains(parse, "--!quest") {
		qData := strings.Split(parse, "|")
		if len(qData) != 3 {
			log.Print("Failed to parse line", parse, "due to split of | being != 3 (", len(qData), ")\n")
			return true
		}
		if err = db.instance.Get(&questData, "SELECT * from quests WHERE tag = ?", qData[1]); err != nil {
			//This is OK, we just insert the data with what we know.

			if result, vErr := db.instance.NamedExec("REPLACE INTO quests (tag) VALUES (:tag)", questData); err != nil {
				log.Fatal("Failed to insert into quests tag", questData.tag, vErr.Error())
			} else {
				questData.questid, _ = result.LastInsertId()
			}
		}

	}

	if strings.Contains(parse, "#!item") || strings.Contains(parse, "--!item") {
		qData := strings.Split(parse, "|")
		if len(qData) != 2 {
			log.Print("Failed to parse line", parse, "due to split of | being != 2 (", len(qData), ")\n")
			return true
		}

		itemField := qData[1]
		if len(itemField) < 1 {
			log.Print("Failed to parse line", parse, "due to len of second element being empty\n")
			return true
		}
		item_id := 0
		if item_id, err = strconv.Atoi(itemField); err != nil {
			log.Print("Failed to parse line", parse, "due to non numberic second element\n")
			return true
		}

		insertQuery := "REPLACE INTO zone_drops (quest_id, item_id, npc_id, zone_id, quest_stepnum, is_quest_item) VALUES (?, ?, ?, ?, ?, ?)"
		stmt, _ := db.instance.Prepare(insertQuery)
		if _, err = stmt.Exec(questData.questid, item_id, npc_id, zone_id, questData.stepnum, 1); err != nil {
			log.Fatal(err.Error())
		}
		stmt.Close()
		return true
	}

	if strings.Contains(parse, "#!reward") || strings.Contains(parse, "--!reward") {
		qData := strings.Split(parse, "|")
		if len(qData) != 2 {
			log.Print("Failed to parse line", parse, "due to split of | being != 2 (", len(qData), ")\n")
			return true
		}

		itemField := qData[1]
		if len(itemField) < 1 {
			log.Print("Failed to parse line", parse, "due to len of second element being empty\n")
			return true
		}
		item_id := 0
		if item_id, err = strconv.Atoi(itemField); err != nil {
			log.Print("Failed to parse line", parse, "due to non numberic second element\n")
			return true
		}

		insertQuery := "REPLACE INTO zone_drops (quest_id, item_id, npc_id, zone_id, quest_stepnum, is_quest_reward) VALUES (?, ?, ?, ?, ?, ?)"
		stmt, _ := db.instance.Prepare(insertQuery)
		if _, err = stmt.Exec(questData.questid, item_id, npc_id, zone_id, questData.stepnum, 1); err != nil {
			log.Fatal(err.Error())
		}
		stmt.Close()
		return true
	}

	if strings.Contains(parse, "#!end") || strings.Contains(parse, "--!end") {
		questData = &QuestData{}
		return true
	}

	return false
}

func getQuestId(shorttag string) {

}

func ignoreParse(parse string) bool {
	return (strings.Contains(parse, "#!ignore") || strings.Contains(parse, "--!ignore"))
}

func endParse(parse string, questData *QuestData) bool {
	if strings.Contains(parse, "#!end") || strings.Contains(parse, "--!end") {
		questData = &QuestData{}
		return true
	}
	return false
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

func showPercent(message string, cur int, max int, insertCount int, color string) {
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
	fmt.Printf("\033[60D")
	val := float64(cur) / float64(max) * float64(dotCount)
	fmt.Printf("%s%s - [", color, message)
	for i := 0; i < dotCount; i++ {
		if int(val) >= i {
			fmt.Printf(".")
		} else {
			fmt.Printf(" ")
		}
	}
	fmt.Printf("]\033[0m")
	if cur == max {
		fmt.Printf(" %d zone entries", insertCount)
		fmt.Printf("\n")
	}
}
