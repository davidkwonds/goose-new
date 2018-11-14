package goose

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"
)

// SQLDialect abstracts the details of specific SQL dialects
// for goose's few SQL specific statements
type SQLDialect interface {
	createVersionTableSQL() string // sql string to create the goose_db_version table
	insertVersionSQL() string      // sql string to insert the initial version table row
	dbVersionQuery(db *sql.DB) (*sql.Rows, error)
	updateTableSQL() string
}

// drivers that we don't know about can ask for a dialect by name
func dialectByName(d string) SQLDialect {
	switch d {
	case "postgres":
		return &PostgresDialect{}
	case "mysql":
		return &MySQLDialect{}
	case "sqlite3":
		return &Sqlite3Dialect{}
	}

	return nil
}

////////////////////////////
// Postgres
////////////////////////////

// PostgresDialect struct
type PostgresDialect struct{}

func (pg PostgresDialect) createVersionTableSQL() string {
	return `CREATE TABLE goose_db_version (
            	id serial NOT NULL,
                version_id bigint NOT NULL,
				is_applied boolean NOT NULL,
				tstamp timestamp NULL default now(),
				work_version varchar(32) NOT NULL,
                PRIMARY KEY(id)
            );`
}

func (pg PostgresDialect) insertVersionSQL() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES ($1, $2);"
}

func (pg PostgresDialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT version_id, is_applied from goose_db_version ORDER BY id DESC")

	// XXX: check for postgres specific error indicating the table doesn't exist.
	// for now, assume any error is because the table doesn't exist,
	// in which case we'll try to create it.
	if err != nil {
		return nil, ErrTableDoesNotExist
	}

	return rows, err
}

func (m *PostgresDialect) updateTableSQL() string {
	return "ALTER TABLE `goose_db_version` ADD `work_version` VARCHAR(32)  NOT NULL"
}

////////////////////////////
// MySQL
////////////////////////////

// MySQLDialect Struct
type MySQLDialect struct{}

func (m MySQLDialect) createVersionTableSQL() string {
	return `CREATE TABLE goose_db_version (
                id serial NOT NULL,
                version_id bigint NOT NULL,
				is_applied boolean NOT NULL,
				tstamp timestamp NULL default now(),
				work_version varchar(32) NOT NULL,
                PRIMARY KEY(id)
            );`
}

func (m MySQLDialect) insertVersionSQL() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES (?, ?);"
}

func (m MySQLDialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT version_id, is_applied from goose_db_version ORDER BY id DESC")

	// XXX: check for mysql specific error indicating the table doesn't exist.
	// for now, assume any error is because the table doesn't exist,
	// in which case we'll try to create it.
	if err != nil {
		return nil, ErrTableDoesNotExist
	}

	return rows, err
}

func (m *MySQLDialect) updateTableSQL() string {
	return "ALTER TABLE `goose_db_version` ADD `work_version` VARCHAR(32)  NOT NULL"
}

////////////////////////////
// sqlite3
////////////////////////////

// Sqlite3Dialect struct
type Sqlite3Dialect struct{}

func (m Sqlite3Dialect) createVersionTableSQL() string {
	return `CREATE TABLE goose_db_version (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                version_id INTEGER NOT NULL,
				is_applied INTEGER NOT NULL,
				tstamp TIMESTAMP DEFAULT (datetime('now'))
				work_version TEXT NOT NULL,
            );`
}

func (m Sqlite3Dialect) insertVersionSQL() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES (?, ?);"
}

func (m Sqlite3Dialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT version_id, is_applied from goose_db_version ORDER BY id DESC")

	switch err.(type) {
	case sqlite3.Error:
		return nil, ErrTableDoesNotExist
	}
	return rows, err
}

func (m *Sqlite3Dialect) updateTableSQL() string {
	return "ALTER TABLE `goose_db_version` ADD `work_version` TEXT  NOT NULL"
}
