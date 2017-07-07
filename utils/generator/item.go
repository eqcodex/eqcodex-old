package main

import (
	//"html/template"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xackery/goeq/item"
	"github.com/xackery/goeq/npc"
)

func generateItem(instance *Instance) {
	var err error

	items := []*item.Item{}

	fmt.Println("Getting items...")
	query := "SELECT * FROM items ORDER by id"
	err = instance.db.Select(&items, query)
	if err != nil {
		log.Println("Failed to select items", err.Error())
		return
	}
	/*if err = os.RemoveAll(instance.yamlConfig.Output + "item/"); err != nil {
		log.Println("Failed to remove item dir contents:", err.Error())
		return
	}*/

	if err = os.MkdirAll(instance.yamlConfig.Output+"item/", 0777); err != nil {
		log.Println("Failed to make item dir:", err.Error())
		return
	}
	max := len(items)
	for focusId, itemEntry := range items {
		if instance.tracker.ItemId > focusId {
			continue
		}
		generateItemEntry(instance, itemEntry)

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
		showPercent(fmt.Sprintf("%d @ %0.2f/sec", focusId, rate), focusId, max, remainString, "green")
		instance.tracker.ItemId = focusId
		saveTracker(instance.tracker)
		//break
		//if zoneEntry.Short_name.String == "airplane" {
		//	break
		//}
	}
	return
}

func generateItemEntry(instance *Instance, itemEntry *item.Item) {
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
		Task_id         int
		Header_line     string
		Class_line      string
		Race_line       string
		Slot_line       string
		Size_line       string
		Weight_line     string
	}

	type NPCData struct {
		*npc.NpcTypes
		Url             string
		Quest           string
		Is_quest_reward int `db:"is_quest_reward"`
		Is_quest_item   int `db:"is_quest_item"`
		Item_id         int
		Npc_id          int
		Zone_name       string
		Zone_url        string
		Zone_id         int
	}

	type PageData struct {
		Core    *CoreData
		Content *ContentData
		Zone    *ZoneData
		Item    *ItemData
		NPCs    []*NPCData
	}

	page := &PageData{
		Core:    getCore(),
		Content: &ContentData{},
		Item: &ItemData{
			Item: itemEntry,
		},
		NPCs: []*NPCData{},
	}

	query := `SELECT npc.*, zone.long_name zone_name, zone.Zoneidnumber zone_id FROM zone_drops
	INNER JOIN zone on zone.zoneidnumber = zone_id
	INNER JOIN npc_types npc ON npc_id = npc.id
	WHERE item_id = ?
	GROUP BY npc.Name`
	if err = instance.db.Select(&page.NPCs, query, itemEntry.Id); err != nil {
		log.Println("Failed to select npcs", err.Error())
		return
	}

	for _, npcEntry := range page.NPCs {
		npcEntry.Quest = ""
		if npcEntry.Is_quest_item == 1 {
			npcEntry.Quest = "Quest Item"
		}
		if npcEntry.Is_quest_reward == 1 {
			npcEntry.Quest = "Quest Reward"
		}

		npcName := cleanUrl(npcEntry.Name)
		if len(npcName) == 0 {
			npcName = "(Blank)"
		}
		npcEntry.Url = fmt.Sprintf("/npc/%s-%d.html", npcName, npcEntry.Id.Int64)
		//log.Println(npcEntry.Name)
		npcEntry.Name = cleanName(npcEntry.Name)
		//log.Println(npcEntry.Name)
		npcEntry.Zone_url = fmt.Sprintf("/zone/%s-%d.html", cleanUrl(npcEntry.Zone_name), npcEntry.Zone_id)
		//log.Println(npcEntry.Zone_url)
	}

	itemUrl := fmt.Sprintf("/item/%s-%d.html", cleanUrl(itemEntry.Name), itemEntry.Id)

	if page.Item.Magic > 0 {
		page.Item.Header_line = "Magic, "
	}
	if page.Item.Notransfer > 0 || page.Item.Nodrop > 0 {
		page.Item.Header_line = "No Trade, "
	}
	if page.Item.Heirloom > 0 {
		page.Item.Header_line = "Heirloom, "
	}
	if len(page.Item.Header_line) > 2 {
		page.Item.Header_line = page.Item.Header_line[0 : len(page.Item.Header_line)-2]
	}

	page.Item.Class_line = getClasses(page.Item.Classes)
	page.Item.Race_line = getRaces(page.Item.Races)
	page.Item.Slot_line = getSlots(page.Item.Slots)
	page.Item.Size_line = getSizes(page.Item.Size)
	page.Item.Weight_line = fmt.Sprintf("%.1f", float64(page.Item.Weight/10))

	page.Core.Site.Title = itemEntry.Name + " | EQCodex"
	t := getCoreTemplate(instance)

	t, err = t.ParseFiles(instance.yamlConfig.Templates + "item/index.tpl")
	if err != nil {
		log.Println("Failed to parse item/index:", err.Error())
		return
	}

	f, err := os.Create(instance.yamlConfig.Output + itemUrl)
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
