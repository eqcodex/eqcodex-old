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

func generateNPC(instance *Instance) {
	var err error
	npcs := []*npc.NpcTypes{}

	fmt.Println("Getting npcs...")
	query := `SELECT npc_types.* FROM npc_types 
	INNER JOIN zone_drops ON npc_id = npc_types.id
	GROUP BY npc_id
	ORDER by npc_id`
	err = instance.db.Select(&npcs, query)
	if err != nil {
		log.Println("Failed to select npcs", err.Error())
		return
	}
	/*if err = os.RemoveAll(instance.yamlConfig.Output + "item/"); err != nil {
		log.Println("Failed to remove item dir contents:", err.Error())
		return
	}*/

	if err = os.MkdirAll(instance.yamlConfig.Output+"npc/", 0777); err != nil {
		log.Println("Failed to make item dir:", err.Error())
		return
	}
	max := len(npcs)
	for focusId, npcEntry := range npcs {
		if instance.tracker.NpcId > focusId {
			continue
		}
		generateNPCEntry(instance, npcEntry)

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
		instance.tracker.NpcId = focusId
		saveTracker(instance.tracker)

	}
	return
}

func generateNPCEntry(instance *Instance, npcEntry *npc.NpcTypes) {
	var err error
	type ContentData struct {
	}

	type ItemData struct {
		*item.Item
		Category        string
		Era             string
		Quest           string
		Url             string
		NPC             string
		Item_id         int
		Npc_id          int
		Zone_id         int
		Is_quest_reward int `db:"is_quest_reward"`
		Is_quest_item   int `db:"is_quest_item"`
		Npc_count       int
		Header_line     string
		Class_line      string
		Race_line       string
		Slot_line       string
		Size_line       string
		Weight_line     string
		Task_id         int
	}

	type NPCData struct {
		*npc.NpcTypes
		Url             string
		Zone_short_name string `db:"short_name"`
		Zone_long_name  string `db:"long_name"`
		Zone_id         int
		Quest           string
		Is_quest_reward int `db:"is_quest_reward"`
		Is_quest_item   int `db:"is_quest_item"`
		MapData         string
	}

	type PageData struct {
		Core    *CoreData
		Content *ContentData
		Zone    *ZoneData
		Npc     *NPCData
		Items   []*ItemData
	}

	page := &PageData{
		Core:    getCore(),
		Content: &ContentData{},
		Npc: &NPCData{
			NpcTypes: npcEntry,
		},
		Items: []*ItemData{},
	}

	npcDump := &NPCData{}
	query := `SELECT zone.short_name, zone.long_name FROM zone_drops 
	INNER JOIN zone ON zone.zoneidnumber = zone_id
	WHERE npc_id = ? LIMIT 1`
	if err = instance.db.Get(npcDump, query, npcEntry.Id.Int64); err != nil {
		//log.Println("Failed to get short_name", err.Error())
		return
	}
	page.Npc.Zone_id = npcDump.Zone_id
	page.Npc.Zone_short_name = npcDump.Zone_short_name
	page.Npc.Zone_long_name = npcDump.Zone_long_name

	page.Npc.MapData, err = mapCanvasByNpcId(npcEntry.Id.Int64, page.Npc.Zone_short_name, instance)
	if err != nil {
		fmt.Println("Failed to parse json", err.Error())
		return
	}
	page.Npc.Name = cleanName(page.Npc.Name)

	query = `SELECT items.* FROM zone_drops
	INNER JOIN items ON items.id = item_id
	WHERE npc_id = ?
	GROUP BY items.name`
	if err = instance.db.Select(&page.Items, query, npcEntry.Id); err != nil {
		log.Println("Failed to select npcs", err.Error())
		return
	}

	for _, itemEntry := range page.Items {
		itemEntry.Quest = ""
		if itemEntry.Is_quest_item == 1 {
			itemEntry.Quest = "Quest Item"
		}
		if itemEntry.Is_quest_reward == 1 {
			itemEntry.Quest = "Quest Reward"
		}

		itemEntry.Category = getCategory(itemEntry.Slots)

		itemEntry.Url = fmt.Sprintf("/item/%s-%d.html", cleanUrl(itemEntry.Name), itemEntry.Id)
		//log.Println(npcEntry.Name)
		//npcEntry.Name = cleanName(npcEntry.Name)
		//log.Println(npcEntry.Name)
		//log.Println(npcEntry.Zone_url)
	}

	page.Core.Site.Title = npcEntry.Name + " | EQCodex"
	t := getCoreTemplate(instance)

	t, err = t.ParseFiles(instance.yamlConfig.Templates + "npc/index.tpl")
	if err != nil {
		log.Println("Failed to parse item/index:", err.Error())
		return
	}

	npcName := cleanUrl(npcEntry.Name)
	if len(npcName) == 0 {
		npcName = "(Blank)"
	}
	npcUrl := fmt.Sprintf("/npc/%s-%d.html", npcName, npcEntry.Id.Int64)

	f, err := os.Create(instance.yamlConfig.Output + npcUrl)
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
