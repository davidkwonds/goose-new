package main

import (
	"fmt"
	"goose/lib/goose"
	"log"
	"os"
)

var mvToWorkVersionCmd = &Command{
	Name:    "mv",
	Usage:   "mv file",
	Summary: "move target file to work version",
	Help:    `mv extended help here...`,
	Run:     mvToWorkVersionRun,
}

func mvToWorkVersionRun(cmd *Command, args ...string) {
	if *flagVersion == "" {
		log.Fatal("goose mv: work version shuould be present")
	}

	if len(args) < 1 {
		log.Fatal("goose mv: file name required")
	}

	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	if err = os.MkdirAll(conf.GetMigrationDir(), 0777); err != nil {
		log.Fatal(err)
	}

	filename := args[0]
	path := fmt.Sprintf("%s/%s", conf.MigrationsDir, filename)
	_, err = os.Stat(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("goose mv: file is not exists, %s", path))
	}

	db, e := goose.OpenDBFromDBConf(conf)
	if e != nil {
		log.Fatal("couldn't open DB:", e)
	}
	defer db.Close()

	if e := goose.MoveFileToWorkVersion(conf, filename, db); e != nil {
		log.Fatal(e)
	}

	fmt.Printf("%s File move Complete\n", path)
}
