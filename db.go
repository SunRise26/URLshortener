package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"time"
)

type database struct {
	*sql.DB
}

func OpenDatabase() *database {
	db, err := sql.Open("sqlite3", "./database/service.db")
	if err != nil {
		log.Fatal(err)
	}
	statement, _ := db.Prepare(
		"CREATE TABLE IF NOT EXISTS links (id INTEGER PRIMARY KEY, origURL TEXT, shortURL TEXT)",
	)
	statement.Exec()
	statement, _ = db.Prepare(
		"CREATE TABLE IF NOT EXISTS requests (tm INTEGER, origURL TEXT, shortURL TEXT)",
	)
	statement.Exec()
	return &database{db}
}

func (db *database) addLink(origURL string, shortURL string) {
	statement, _ := db.Prepare("INSERT INTO links (origURL, shortURL) VALUES (?, ?)")

	statement.Exec(origURL, shortURL)
}

func (db *database) deleteLink(shortURL string) {
	statement, _ := db.Prepare("DELETE FROM links WHERE shortURL = (?)")

	statement.Exec(shortURL)
}

func (db *database) addRequest(origURL string, shortURL string) {
	tm := time.Now().Unix()
	statement, _ := db.Prepare("INSERT INTO requests (tm, origURL, shortURL) VALUES (?, ?, ?)")

	statement.Exec(tm, origURL, shortURL)
}

func (db *database) deleteOldRequests(lifetime int64) {
	statement, _ := db.Prepare("DELETE FROM requests WHERE ((?) - tm) > (?)")

	statement.Exec(time.Now().Unix(), lifetime)
}

func (db *database) getLinks() string {
	var id	int
	var table, origURL, shortURL string

	rows, _ := db.Query("SELECT id, origURL, shortURL FROM links")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id, &origURL, &shortURL)
		table += strconv.Itoa(id) + ": " + origURL + "\t" + shortURL + "\n"
	}
	return table
}

func (db *database) getRequests() string {
	var tm int64
	var table, origURL, shortURL string

	rows, _ := db.Query("SELECT tm, origURL, shortURL FROM requests")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&tm, &origURL, &shortURL)
		table += time.Unix(tm, 0).String() + ": " + origURL + "\t" + shortURL + "\n"
	}
	return table
}

func (db *database) getOrigURL (shortURL string) string {
	var origURL string

	rows, _ := db.Query("SELECT origURL FROM links WHERE shortURL = (?)", shortURL)
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&origURL)
	}
	return origURL
}
