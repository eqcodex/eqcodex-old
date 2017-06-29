package main

import (
	//"html/template"
	"fmt"
	"log"
	"os"

	"github.com/xackery/goeq/item"
	"github.com/xackery/goeq/npc"
)

func generateItem(instance *Instance) {
	var err error

	items := []*item.Item{}

	fmt.Println("Getting items...")
	query := "SELECT * FROM items ORDER by id" // LIMIT 1"
	err = instance.db.Select(&items, query)
	if err != nil {
		log.Println("Failed to select zones", err.Error())
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
	itemCount := 0
	for _, itemEntry := range items {
		itemCount++
		generateItemEntry(instance, itemEntry)
		if itemCount%1000 == 0 {
			fmt.Printf("%d (%.2f)", itemCount, float64(itemCount/len(items)))
		}
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

	query := `SELECT npc.*, zone.long_name zone_name FROM zone_drops
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

		npcEntry.Url = fmt.Sprintf("/npc/%s-%d.html", cleanUrl(npcEntry.Name), npcEntry.Id.Int64)
		//log.Println(npcEntry.Name)
		npcEntry.Name = cleanName(npcEntry.Name)
		//log.Println(npcEntry.Name)
		npcEntry.Zone_url = fmt.Sprintf("/zone/%s.html", cleanUrl(npcEntry.Zone_name))
		//log.Println(npcEntry.Zone_url)
	}

	itemUrl := fmt.Sprintf("/item/%s-%d.html", cleanUrl(itemEntry.Name), itemEntry.Id)

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
