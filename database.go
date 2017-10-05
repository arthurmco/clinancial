package main

/*
 *  Define database connection things
 *  Copyright (C) 2017 Arthur Mendes
 */
import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var dbpath string = "clinancial.db"

func SetDatabasePath(s string) {
	dbpath = s
}

func GetDatabasePath() string {
	return dbpath
}

func CreateDatabase() error {
	db, eerr := sql.Open("sqlite3", GetDatabasePath())
	if eerr != nil {
		return eerr
	}

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS accounts (" +
		"id INTEGER PRIMARY KEY, name TEXT, ctime INTEGER)")
	if err != nil {
		return err
	}
	stmt.Exec()

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS registers (" +
		"id INTEGER PRIMARY KEY, sid INTEGER, name string, " +
		"time INTEGER, val REAL, fromaccount INTEGER, " +
		"toaccount INTEGER) ")
	if err != nil {
		return err
	}
	stmt.Exec()
	db.Close()
	return nil
}

func DropDatabase() error {
	db, eerr := sql.Open("sqlite3", GetDatabasePath())
	if eerr != nil {
		return eerr
	}

	stmt, err := db.Prepare("DROP TABLE IF EXISTS accounts")
	if err != nil {
		return err
	}
	stmt.Exec()

	stmt, err = db.Prepare("DROP TABLE IF EXISTS registers")
	if err != nil {
		return err
	}
	stmt.Exec()
	db.Close()
	return nil
}
