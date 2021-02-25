package main

import (
	"database/sql"
	"log"
	"time"

)

//Db represents a database object
var Db *sql.DB

//Bug is the in-app representation of a bug
type Bug struct {
	user         *User
	BugID       string
	Title       string
	Description string
	CreatedAt   time.Time
	AssignedTo  Dev
}

//Dev represents a developer under a company
type Dev struct {
	DevName  string
	DevID    string
	role     string
	stack    string
	Email    string
	Password string
	user      *User
	issues   []Bug
}





func init() {
	loadConfig()
	Db, err := sql.Open(config.driverName, config.datasource)
	if err != nil {
		log.Fatal(err)
	}
	Db.Close()
}
