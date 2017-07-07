package main

import (
	//"html/template"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xackery/goeq/item"
	"github.com/xackery/goeq/zone"
)

type ZoneData struct {
	*zone.Zone
	Url string
}

func (i *ZoneData) IsLevel(levels int, val int) bool {
	switch val {
	case 1:
		val = 1
	case 5:
		val = 2
	case 10:
		val = 4
	case 15:
		val = 8
	case 20:
		val = 16
	case 25:
		val = 32
	case 30:
		val = 64
	case 35:
		val = 128
	case 40:
		val = 256
	case 45:
		val = 512
	case 50:
		val = 1024
	case 55:
		val = 2048
	case 60:
		val = 4096
	case 65:
		val = 8192
	case 70:
		val = 16384
	}
	return levels&val == val
}

func generateZoneList(instance *Instance) {
	var err error

	type ContentData struct {
	}

	type PageData struct {
		Core    *CoreData
		Content *ContentData
		Zones   []*ZoneData
	}

	page := &PageData{
		Core:    getCore(),
		Content: &ContentData{},
	}

	//Get zones.
	query := "SELECT * FROM zone ORDER BY short_name"
	err = instance.db.Select(&page.Zones, query)
	if err != nil {
		log.Println("Failed to select zones", err.Error())
		return
	}

	if err = os.MkdirAll(instance.yamlConfig.Output+"zone/", 0777); err != nil {
		log.Println("Failed to make zone dir:", err.Error())
		return
	}

	max := len(page.Zones)
	for focusId, zoneEntry := range page.Zones {

		zoneEntry.Url = fmt.Sprintf("zone/%s-%d.html", cleanUrl(zoneEntry.Long_name.String), zoneEntry.Id.Int64)
		generateZoneEntry(instance, zoneEntry)
		rate := float64(focusId) / float64(time.Since(startTime).Seconds())
		remainString := ""
		remain := (float64(max) - float64(focusId)) / rate

		if remain > 60 {
			remain = remain / 60
			if remain > 60 {
				remain = remain / 60
				remainString = fmt.Sprintf("%0.2f hours", remain)
			} else {
				remainString = fmt.Sprintf("%0.1f minutes", remain)
			}
		} else {
			remainString = fmt.Sprintf("%0.0f seconds", remain)
		}
		showPercent(fmt.Sprintf("%s %d @ %0.2f/sec", zoneEntry.Short_name.String, focusId, rate), focusId, max, remainString, "green")
		//if zoneEntry.Short_name.String == "airplane" {
		//	break
		//}
	}
	page.Core.Site.Title = "Zone Leveling Chart | EQCodex"
	t := getCoreTemplate(instance)

	t, err = t.ParseFiles(instance.yamlConfig.Templates + "zone.tpl")
	if err != nil {
		log.Println("Failed to parse zone:", err.Error())
		return
	}

	f, err := os.Create(instance.yamlConfig.Output + "zone.html")
	if err != nil {
		log.Println("Failed to create file:", err.Error())
		return
	}

	err = t.Execute(f, page)
	if err != nil {
		log.Println("Failed to execute index:", err.Error())
		return
	}
	f.Close()

	return
}

func generateZoneEntry(instance *Instance, zoneEntry *ZoneData) {
	var err error
	type ContentData struct {
	}

	type ItemData struct {
		*item.Item
		Category        string
		Era             string
		Quest           string
		NPC             string
		Item_id         int
		Npc_id          int
		Zone_id         int
		Is_quest_reward int `db:"is_quest_reward"`
		Is_quest_item   int `db:"is_quest_item"`
		Npc_count       int
		Url             string
		Task_id         int
	}

	type PageData struct {
		Core    *CoreData
		Content *ContentData
		Zone    *ZoneData
		Items   []*ItemData
	}

	page := &PageData{
		Core:    getCore(),
		Content: &ContentData{},
		Zone:    zoneEntry,
	}

	query := `SELECT * FROM zone_drops 
	INNER JOIN items ON items.id = zone_drops.item_id
	WHERE zone_id = ? 
	GROUP BY item_id`
	if err = instance.db.Select(&page.Items, query, zoneEntry.Zoneidnumber); err != nil {
		log.Println("Failed to select zone_drops", err.Error())
		return
	}

	for _, itemEntry := range page.Items {
		itemEntry.Category = getCategory(itemEntry.Slots)
		itemEntry.Era = fmt.Sprintf("%d", zoneEntry.Expansion)
		itemEntry.Quest = ""
		if itemEntry.Is_quest_item == 1 {
			itemEntry.Quest = "Quest Item"
		}
		if itemEntry.Is_quest_reward == 1 {
			itemEntry.Quest = "Quest Reward"
		}
		query = `SELECT npc.name FROM zone_drops
		INNER JOIN npc_types npc ON npc.id = zone_drops.npc_id
		WHERE item_id = ? 
		LIMIT 1`

		if err = instance.db.Get(&itemEntry.NPC, query, itemEntry.Id); err != nil {
			log.Println("Failed to select item", err.Error())
			//return
		}

		query = `SELECT count(npc_id) npc_count FROM zone_drops WHERE item_id = ?`
		if err = instance.db.Get(&itemEntry.Npc_count, query, itemEntry.Id); err != nil {
			log.Println("Failed to get count of npcs for item", err.Error())
		}
		itemEntry.Npc_count--
		itemEntry.NPC = cleanName(itemEntry.NPC)
		itemEntry.Url = fmt.Sprintf("/item/%s-%d.html", cleanUrl(itemEntry.Name), itemEntry.Id)
		if itemEntry.Npc_count > 1 {
			itemEntry.NPC += fmt.Sprintf(" and %d more NPCs", itemEntry.Npc_count)
		}
	}

	page.Core.Site.Title = zoneEntry.Long_name.String + " | EQCodex"
	t := getCoreTemplate(instance)

	t, err = t.ParseFiles(instance.yamlConfig.Templates + "zone/index.tpl")
	if err != nil {
		log.Println("Failed to parse zone/index:", err.Error())
		return
	}

	f, err := os.Create(instance.yamlConfig.Output + zoneEntry.Url)
	if err != nil {
		log.Println("Failed to create file:", err.Error())
		return
	}

	err = t.Execute(f, page)
	if err != nil {
		log.Println("Failed to execute index:", err.Error())
		return
	}
	f.Close()
}
