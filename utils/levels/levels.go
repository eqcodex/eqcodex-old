//Generates levels entry of zone
package main

import ( //"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xackery/eqemuconfig"
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
	startTime := time.Now()
	if err = loadConfig(); err != nil {
		log.Fatal(err.Error())
	}

	db := Database{}
	if err = loadDatabase(&db); err != nil {
		log.Fatal(err.Error())
	}

	zones := []zone.Zone{}
	if zones, err = db.getZones(); err != nil {
		log.Fatal(err.Error())
	}

	//iterate zones
	for _, zone := range zones {
		fmt.Print("\n" + zone.Short_name.String + "...")
		spawns := []spawn.Spawn2{}
		if spawns, err = db.getSpawns(zone.Short_name.String); err != nil {
			log.Fatal(err.Error())
		}
		if zone.Min_status > 80 {
			fmt.Printf("Skipping, status > 0, %d", zone.Min_status)
			continue
		}
		mobCounter := map[int]int{}

		//get spawns
		for _, spawn2 := range spawns {
			spawnentries := []spawn.SpawnEntry{}
			if spawnentries, err = db.getSpawnEntries(spawn2.Spawngroupid); err != nil {
				log.Fatal("spawnentry:", err.Error())
			}
			//get npcs
			for _, spawnentry := range spawnentries {
				npcs := []npc.NpcTypes{}
				if npcs, err = db.getNpcs(spawnentry.Npcid); err != nil {
					log.Fatal("npcs:", err.Error())
				}
				for _, npc := range npcs {
					//can't kill non-class mobs
					if npc.Class < 1 && npc.Class > 17 {
						continue
					}
					if strings.Contains(npc.Special_abilities.String, "1^35") ||
						strings.Contains(npc.Special_abilities.String, "1^25") ||
						strings.Contains(npc.Special_abilities.String, "1^24") { //invul
						continue
					}

					if npc.Level == 1 { //1
						mobCounter[1] += 1
					}
					if npc.Level < 6 && npc.Level > 2 { //5
						mobCounter[2] += 1
					}
					if npc.Level < 11 && npc.Level > 7 { //10
						mobCounter[4] += 1
					}
					if npc.Level < 16 && npc.Level > 11 { //15
						mobCounter[8] += 1
					}
					if npc.Level < 21 && npc.Level > 15 { //20
						mobCounter[16] += 1
					}
					if npc.Level < 26 && npc.Level > 20 { //25
						mobCounter[32] += 1
					}
					if npc.Level < 31 && npc.Level > 25 { //30
						mobCounter[64] += 1
					}
					if npc.Level < 36 && npc.Level > 25 { //35
						mobCounter[128] += 1
					}
					if npc.Level < 41 && npc.Level > 35 { //40
						mobCounter[256] += 1
					}
					if npc.Level < 46 && npc.Level > 40 { //45
						mobCounter[512] += 1
					}
					if npc.Level < 51 && npc.Level > 45 { //50
						mobCounter[1024] += 1
					}
					if npc.Level < 56 && npc.Level > 50 { //55
						mobCounter[2048] += 1
					}
					if npc.Level < 61 && npc.Level > 55 { //60
						mobCounter[4096] += 1
					}
					if npc.Level < 66 && npc.Level > 60 { //65
						mobCounter[8192] += 1
					}
					if npc.Level < 71 && npc.Level > 65 { //70
						mobCounter[16384] += 1
					}
				}

			}

		}

		zone.Levels = 0
		minCount := 20
		if mobCounter[1] >= minCount {
			zone.Levels |= 1
		}
		if mobCounter[2] >= minCount {
			zone.Levels |= 2
		}
		if mobCounter[4] >= minCount {
			zone.Levels |= 4
		}
		if mobCounter[8] >= minCount {
			zone.Levels |= 8
		}
		if mobCounter[16] >= minCount {
			zone.Levels |= 16
		}
		if mobCounter[32] >= minCount {
			zone.Levels |= 32
		}
		if mobCounter[64] >= minCount {
			zone.Levels |= 64
		}
		if mobCounter[128] >= minCount {
			zone.Levels |= 128
		}
		if mobCounter[256] >= minCount {
			zone.Levels |= 256
		}
		if mobCounter[512] >= minCount {
			zone.Levels |= 512
		}
		if mobCounter[1024] >= minCount {
			zone.Levels |= 1024
		}
		if mobCounter[2048] >= minCount {
			zone.Levels |= 2048
		}
		if mobCounter[4096] >= minCount {
			zone.Levels |= 4096
			if zone.Levels&1 == 1 {
				zone.Levels -= 1
			}
			if zone.Levels&2 == 2 {
				zone.Levels -= 2
			}
			if zone.Levels&4 == 4 {
				zone.Levels -= 4
			}
		}
		if mobCounter[8192] >= minCount {
			zone.Levels |= 8192
			if zone.Levels&1 == 1 {
				zone.Levels -= 1
			}
			if zone.Levels&2 == 2 {
				zone.Levels -= 2
			}
			if zone.Levels&4 == 4 {
				zone.Levels -= 4
			}
		}
		if mobCounter[16384] >= minCount {
			zone.Levels |= 16384
			if zone.Levels&1 == 1 {
				zone.Levels -= 1
			}
			if zone.Levels&2 == 2 {
				zone.Levels -= 2
			}
			if zone.Levels&4 == 4 {
				zone.Levels -= 4
			}
		}
		query := "UPDATE zone SET levels = ? WHERE zoneidnumber = ?"
		_, err = db.instance.Exec(query, zone.Levels, zone.Zoneidnumber)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%d", zone.Levels)
		//fmt.Printf("(%d)", (insertCount - lastInsertCount))
	}
	log.Println("Done in", time.Since(startTime).Seconds(), "seconds")
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
