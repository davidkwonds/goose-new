package main

import (
	"log"

	"github.com/davidkwonds/goose-new/goose/lib/goose"
)

var upCmd = &Command{
	Name:    "up",
	Usage:   "",
	Summary: "Migrate the DB to the most recent version available",
	Help:    `up extended help here...`,
	Run:     upRun,
}

func upRun(cmd *Command, args ...string) {

	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	target, err := goose.GetMostRecentDBVersion(conf.GetMigrationDir())
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.RunMigrations(conf, target); err != nil {
		log.Fatal(err)
	}
}
