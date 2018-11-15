package goose

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/lib/pq"
)

// DBDriver encapsulates the info needed to work with
// a specific database driver
type DBDriver struct {
	Name    string
	OpenStr string
	Import  string
	Dialect SQLDialect
}

// DBConf struct
type DBConf struct {
	MigrationsDir string
	Env           string
	Driver        DBDriver
	PgSchema      string
	WorkVersion   string
}

// GetMigrationDir get mirgration dir
func (m *DBConf) GetMigrationDir() string {
	return fmt.Sprintf("%s/%s", m.MigrationsDir, m.WorkVersion)
}

// NewDBConf extract configuration details from the given file
func NewDBConf(p, database, env, pgschema, workVersion string) (*DBConf, error) {

	cfgFile := filepath.Join(p, "dbconf.yml")

	f, err := yaml.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	drv, err := f.Get(fmt.Sprintf("%s.%s.driver", database, env))
	if err != nil {
		return nil, err
	}
	drv = os.ExpandEnv(drv)

	open, err := f.Get(fmt.Sprintf("%s.%s.open", database, env))
	if err != nil {
		return nil, err
	}
	open = os.ExpandEnv(open)

	// Automatically parse postgres urls
	if drv == "postgres" {

		// Assumption: If we can parse the URL, we should
		if parsedURL, err := pq.ParseURL(open); err == nil && parsedURL != "" {
			open = parsedURL
		}
	}

	d := newDBDriver(drv, open)

	// allow the configuration to override the Import for this driver
	if imprt, err := f.Get(fmt.Sprintf("%s.%s.import", database, env)); err == nil {
		d.Import = imprt
	}

	// allow the configuration to override the Dialect for this driver
	if dialect, err := f.Get(fmt.Sprintf("%s.%s.dialect", database, env)); err == nil {
		d.Dialect = dialectByName(dialect)
	}

	if !d.IsValid() {
		return nil, fmt.Errorf("Invalid DBConf: %v", d)
	}

	return &DBConf{
		MigrationsDir: filepath.Join(p, database, "migrations"),
		Env:           env,
		Driver:        d,
		PgSchema:      pgschema,
		WorkVersion:   workVersion,
	}, nil
}

// Create a new DBDriver and populate driver specific
// fields for drivers that we know about.
// Further customization may be done in NewDBConf
func newDBDriver(name, open string) DBDriver {

	d := DBDriver{
		Name:    name,
		OpenStr: open,
	}

	switch name {
	case "postgres":
		d.Import = "github.com/lib/pq"
		d.Dialect = &PostgresDialect{}

	case "mymysql":
		d.Import = "github.com/ziutek/mymysql/godrv"
		d.Dialect = &MySQLDialect{}

	case "mysql":
		d.Import = "github.com/go-sql-driver/mysql"
		d.Dialect = &MySQLDialect{}

	case "sqlite3":
		d.Import = "github.com/mattn/go-sqlite3"
		d.Dialect = &Sqlite3Dialect{}
	}

	return d
}

// IsValid ensure we have enough info about this driver
func (drv *DBDriver) IsValid() bool {
	return len(drv.Import) > 0 && drv.Dialect != nil
}

// OpenDBFromDBConf wraps database/sql.DB.Open() and configures
// the newly opened DB based on the given DBConf.
//
// Callers must Close() the returned DB.
func OpenDBFromDBConf(conf *DBConf) (*sql.DB, error) {
	db, err := sql.Open(conf.Driver.Name, conf.Driver.OpenStr)
	if err != nil {
		return nil, err
	}

	// if a postgres schema has been specified, apply it
	if conf.Driver.Name == "postgres" && conf.PgSchema != "" {
		if _, err := db.Exec("SET search_path TO " + conf.PgSchema); err != nil {
			return nil, err
		}
	}

	return db, nil
}
