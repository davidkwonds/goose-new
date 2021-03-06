package main

import (
	"log"

	"github.com/davidkwonds/goose-new/lib/goose"
)

var downCmd = &Command{
	Name:    "down",
	Usage:   "",
	Summary: "Roll back the version by 1",
	Help:    `down extended help here...`,
	Run:     downRun,
}

func downRun(cmd *Command, args ...string) {

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

	if err = goose.RunMigrations(conf, previous); err != nil {
		log.Fatal(err)
	}
}
