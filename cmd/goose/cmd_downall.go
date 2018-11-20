package main

import (
	"log"

	"github.com/davidkwonds/goose-new/lib/goose"
)

var downAllCmd = &Command{
	Name:    "downall",
	Usage:   "",
	Summary: "Roll back all in the work version",
	Help:    `down extended help here...`,
	Run:     downAllRun,
}

func downAllRun(cmd *Command, args ...string) {

	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	for {
		current, err := goose.GetDBVersion(conf)
		if err != nil {
			log.Fatal(err)
		}
		if current == 0 {
			break
		}

		previous, err := goose.GetPreviousDBVersion(conf.GetMigrationDir(), current)
		if err != nil {
			log.Fatal(err)
		}

		if err = goose.RunMigrations(conf, previous); err != nil {
			log.Fatal(err)
		}
	}
}
