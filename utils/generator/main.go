//This program generates HTML files
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xackery/eqemuconfig"
)

type Instance struct {
	yamlConfig  *YamlConfig
	eqemuConfig *eqemuconfig.Config
	db          *sqlx.DB
}

func main() {
	instance := &Instance{}
	done := make(chan bool)
	watchForChanges(instance)
	go serveWebsite(instance)
	<-done
}

func serveWebsite(instance *Instance) {
	fs := http.FileServer(http.Dir(instance.yamlConfig.Www))
	http.Handle("/", fs)
	log.Println("Listening on :3000...")
	http.ListenAndServe(":3000", nil)
}

func generateTemplates(instance *Instance) {
	var err error
	startTime := time.Now()
	fmt.Println("Staring up...")

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
	generateZoneList(instance)
	generateItem(instance)
	generateIndex(instance)
	fmt.Println("Completed in", time.Since(startTime).Seconds())
}
