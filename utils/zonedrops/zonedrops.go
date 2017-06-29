//Zonedrops parses zones and creates a pivot for item drops
package main

/*
#iterate zones
#SELECT short_name, zoneidnumber, long_name FROM zone;
#get spawngroupid of zone
#SELECT spawngroupid FROM spawn2 WHERE zone = "crushbone";
#get npcs of spawngroup
#SELECT npcid FROM spawnentry WHERE spawngroupid = 541;
#get loottable_id of npc
#SELECT NAME, loottable_id FROM npc_types WHERE id = 58007;
#get lootdrops
#SELECT lootdrop_id FROM loottable_entries WHERE loottable_id = 1719;
#get itemid of drops
#SELECT item_id FROM lootdrop_entries WHERE lootdrop_id = 3348;
#get item data.
#SELECT * FROM items WHERE id = 5040;
*/
import ( //"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/goeq/item"
	"github.com/xackery/goeq/loot"
	"github.com/xackery/goeq/npc"
	"github.com/xackery/goeq/spawn"
	"github.com/xackery/goeq/zone"
)

var (
	config *eqemuconfig.Config
)

type Database struct {
	instance *sqlx.DB
}

func main() {
	var err error
	if err = loadConfig(); err != nil {
		log.Fatal(err.Error())
	}

	db := Database{}
	if err = loadDatabase(&db); err != nil {
		log.Fatal(err.Error())
	}

	if err = db.createTable(); err != nil {
		log.Fatal(err.Error())
	}
	//if err = db.truncateTable(); err != nil {
	//	log.Fatal(err.Error())
	//}

	insertCount := 0
	zones := []zone.Zone{}
	if zones, err = db.getZones(); err != nil {
		log.Fatal(err.Error())
	}

	insertQuery := "REPLACE INTO zone_drops (item_id, npc_id, zone_id) VALUES"
	insertVals := []interface{}{}

	isSkip := true

	for _, zone := range zones {
		//skipping feature
		if zone.Short_name.String != "steamfactory" && isSkip {
			fmt.Println("Skipping", zone.Short_name.String)
			continue
		}
		isSkip = false
		lastInsertCount := insertCount
		fmt.Print("\n" + zone.Short_name.String + "...")
		spawns := []spawn.Spawn2{}
		if spawns, err = db.getSpawns(zone.Short_name.String); err != nil {
			log.Fatal(err.Error())
		}
		if zone.Min_status > 80 {
			fmt.Printf("Skipping, status > 0, %d", zone.Min_status)
			continue
		}

		for _, spawn2 := range spawns {
			spawnentries := []spawn.SpawnEntry{}
			if spawnentries, err = db.getSpawnEntries(spawn2.Spawngroupid); err != nil {
				log.Fatal("spawnentry:", err.Error())
			}
			for _, spawnentry := range spawnentries {
				npcs := []npc.NpcTypes{}
				if npcs, err = db.getNpcs(spawnentry.Npcid); err != nil {
					log.Fatal("npcs:", err.Error())
				}
				for _, npc := range npcs {
					loottableentries := []loot.LootTableEntries{}
					if loottableentries, err = db.getLoottableEntries(npc.Loottable_id); err != nil {
						log.Fatal("loottable:", err.Error())
					}
					for _, loottableentry := range loottableentries {
						lootdropentries := []loot.LootDropEntries{}
						if lootdropentries, err = db.getLootdropEntries(loottableentry.Lootdrop_id); err != nil {
							log.Fatal("lootdrop:", err.Error())
						}
						for _, lootdropentry := range lootdropentries {
							items := []item.Item{}
							if items, err = db.getItems(lootdropentry.Item_id); err != nil {
								log.Fatal("items", err.Error())
							}
							for _, itementry := range items {
								insertCount++
								if insertCount%1000 == 0 {
									fmt.Printf("%d, ", insertCount)
								}
								insertQuery += "(?, ?, ?),"
								insertVals = append(insertVals, itementry.Id, npc.Id.Int64, zone.Zoneidnumber)

								if insertCount%5000 == 0 {
									fmt.Print(", inserting 5k records...")
									insertQuery = insertQuery[0 : len(insertQuery)-1]
									stmt, _ := db.instance.Prepare(insertQuery)

									if _, err = stmt.Exec(insertVals...); err != nil {
										log.Fatal(err.Error())
									}
									//reset query
									insertQuery = "REPLACE INTO zone_drops (item_id, npc_id, zone_id) VALUES"
									insertVals = []interface{}{}
								}
								//fmt.Printf(itementry.Name + ", ")
								//db.instance.Exec("(?,?,?,?)", itementry.Id, npc.Id, zone.Short_name, zone.Id)
							}
						}
					}
				}
			}

		}
		fmt.Printf("(%d)", (insertCount - lastInsertCount))
	}
	//trim the last
	insertQuery = insertQuery[0 : len(insertQuery)-1]
	stmt, _ := db.instance.Prepare(insertQuery)
	if _, err = stmt.Exec(insertVals); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Done!", insertCount)
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

func (db *Database) createTable() error {
	if db == nil {
		return fmt.Errorf("No database")
	}
	if db.instance == nil {
		return fmt.Errorf("No database instance")
	}
	query := "CREATE TABLE IF NOT EXISTS `zone_drops` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`item_id` int(11) DEFAULT NULL," +
		"`npc_id` int(11) DEFAULT NULL," +
		"`zone_short_name` int(11) DEFAULT NULL," +
		"`zone_id` int(11) DEFAULT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `item_id` (`item_id`,`npc_id`,`zone_id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	if _, err := db.instance.Exec(query); err != nil {
		return err
	}
	return nil
}

func (db *Database) truncateTable() error {
	if db == nil {
		return fmt.Errorf("No database")
	}
	if db.instance == nil {
		return fmt.Errorf("No database instance")
	}
	query := "TRUNCATE zone_drops"
	if _, err := db.instance.Exec(query); err != nil {
		return err
	}
	return nil
}

func (db *Database) getZones() ([]zone.Zone, error) {
	if db == nil {
		return nil, fmt.Errorf("No database")
	}
	if db.instance == nil {
		return nil, fmt.Errorf("No database instance")
	}
	zones := []zone.Zone{}
	query := "SELECT * from eqemu.zone WHERE min_status < 90 ORDER BY short_name ASC"
	if err := db.instance.Select(&zones, query); err != nil {
		return nil, err
	}
	return zones, nil
}

func (db *Database) getSpawns(zoneShortname string) ([]spawn.Spawn2, error) {
	if db == nil {
		return nil, fmt.Errorf("No database")
	}
	if db.instance == nil {
		return nil, fmt.Errorf("No database instance")
	}
	spawns := []spawn.Spawn2{}
	query := "SELECT * FROM spawn2 WHERE zone = ?"
	if err := db.instance.Select(&spawns, query, zoneShortname); err != nil {
		return nil, err
	}
	return spawns, nil
}

func (db *Database) getSpawnEntries(spawngroupid int) ([]spawn.SpawnEntry, error) {
	if db == nil {
		return nil, fmt.Errorf("No database")
	}
	if db.instance == nil {
		return nil, fmt.Errorf("No database instance")
	}
	spawnentries := []spawn.SpawnEntry{}
	query := "SELECT * FROM spawnentry WHERE spawngroupid = ?"
	if err := db.instance.Select(&spawnentries, query, spawngroupid); err != nil {
		return nil, err
	}
	return spawnentries, nil
}

func (db *Database) getNpcs(npctypeid int) ([]npc.NpcTypes, error) {
	if db == nil {
		return nil, fmt.Errorf("No database")
	}
	if db.instance == nil {
		return nil, fmt.Errorf("No database instance")
	}
	npcs := []npc.NpcTypes{}
	query := "SELECT * FROM npc_types WHERE id = ?"
	if err := db.instance.Select(&npcs, query, npctypeid); err != nil {
		return nil, err
	}
	return npcs, nil
}

func (db *Database) getLoottableEntries(loottableid int) ([]loot.LootTableEntries, error) {
	if db == nil {
		return nil, fmt.Errorf("No database")
	}
	if db.instance == nil {
		return nil, fmt.Errorf("No database instance")
	}
	loottableentries := []loot.LootTableEntries{}
	query := "SELECT * FROM loottable_entries WHERE loottable_id = ?"
	if err := db.instance.Select(&loottableentries, query, loottableid); err != nil {
		return nil, err
	}
	return loottableentries, nil
}

func (db *Database) getLootdropEntries(lootdropid int) ([]loot.LootDropEntries, error) {
	if db == nil {
		return nil, fmt.Errorf("No database")
	}
	if db.instance == nil {
		return nil, fmt.Errorf("No database instance")
	}
	lootdropentries := []loot.LootDropEntries{}
	query := "SELECT * FROM lootdrop_entries WHERE lootdrop_id = ?"
	if err := db.instance.Select(&lootdropentries, query, lootdropid); err != nil {
		return nil, err
	}
	return lootdropentries, nil
}

func (db *Database) getItems(itemid int) ([]item.Item, error) {
	if db == nil {
		return nil, fmt.Errorf("No database")
	}
	if db.instance == nil {
		return nil, fmt.Errorf("No database instance")
	}
	items := []item.Item{}
	query := "SELECT * FROM items WHERE id = ?"
	if err := db.instance.Select(&items, query, itemid); err != nil {
		return nil, err
	}
	return items, nil
}
