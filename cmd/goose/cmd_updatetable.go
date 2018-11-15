package main

import (
	"fmt"
	"log"

	"github.com/davidkwonds/goose-new/goose/lib/goose"
)

var updatetableCmd = &Command{
	Name:    "updatetable",
	Usage:   "",
	Summary: "update old goose table to new",
	Help:    `dbversion extended help here...`,
	Run:     updatetableRun,
}

func updatetableRun(cmd *Command, args ...string) {
	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	db, e := goose.OpenDBFromDBConf(conf)
	if e != nil {
		log.Fatal("couldn't open DB:", e)
	}
	defer db.Close()

	if e := goose.UpdateOldTable(conf, db); e != nil {
		log.Fatal(e)
	}

	fmt.Printf("%s Table Update Complete", conf.MigrationsDir)
}
