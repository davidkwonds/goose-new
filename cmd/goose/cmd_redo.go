package main

import (
	"goose/lib/goose"
	"log"
)

var redoCmd = &Command{
	Name:    "redo",
	Usage:   "",
	Summary: "Re-run the latest migration",
	Help:    `redo extended help here...`,
	Run:     redoRun,
}

func redoRun(cmd *Command, args ...string) {
	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	current, err := goose.GetDBVersion(conf)
	if err != nil {
		log.Fatal(err)
	}

	previous, err := goose.GetPreviousDBVersion(conf.GetMigrationDir(), current)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.RunMigrations(conf, previous); err != nil {
		log.Fatal(err)
	}

	if err := goose.RunMigrations(conf, current); err != nil {
		log.Fatal(err)
	}
}
