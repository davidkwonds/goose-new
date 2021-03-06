package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/davidkwonds/goose-new/lib/goose"
)

// global options. available to any subcommands.
var flagPath = flag.String("path", "database", "folder containing db info")
var flagDatabase = flag.String("db", "master", "which database to use")
var flagEnv = flag.String("env", "development", "which DB environment to use")
var flagPgSchema = flag.String("pgschema", "", "which postgres-schema to migrate (default = none)")
var flagVersion = flag.String("v", "", "target version")

// helper to create a DBConf from the given flags
func dbConfFromFlags() (dbconf *goose.DBConf, err error) {
	return goose.NewDBConf(*flagPath, *flagDatabase, *flagEnv, *flagPgSchema, *flagVersion)
}

var commands = []*Command{
	upCmd,
	downCmd,
	downAllCmd,
	redoCmd,
	statusCmd,
	createCmd,
	dbVersionCmd,
	updatetableCmd,
	mvToWorkVersionCmd,
}

func main() {

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 || args[0] == "-h" {
		flag.Usage()
		return
	}

	var cmd *Command
	name := args[0]
	for _, c := range commands {
		if strings.HasPrefix(c.Name, name) {
			cmd = c
			break
		}
	}

	if cmd == nil {
		fmt.Printf("error: unknown command %q\n", name)
		flag.Usage()
		os.Exit(1)
	}

	cmd.Exec(args[1:])
}

func usage() {
	fmt.Print(usagePrefix)
	flag.PrintDefaults()
	usageTmpl.Execute(os.Stdout, commands)
}

var usagePrefix = `
goose is a database migration management system for Go projects.

Usage:
    goose [options] <subcommand> [subcommand options]

Options:
`
var usageTmpl = template.Must(template.New("usage").Parse(
	`
Commands:{{range .}}
    {{.Name | printf "%-10s"}} {{.Summary}}{{end}}
`))
