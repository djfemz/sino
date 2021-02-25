package main

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//Configuration represents the preloaded application configuration
type Configuration struct {
	driverName string
	datasource string
	Address    string
}

var config Configuration

func loadConfig() {

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

}

//CreateBugID creates a unique user identification number.
func CreateBugID() (BugID string) {
	var idSlice, row []string

	statement := "select bugid value $1 from bugs"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	err = rows.Scan(row)
	if err != nil {
		log.Fatal(err)
	}
	if row[len(row)-1] != "" {
		idSlice = strings.Split(row[len(row)-1], "-")
		idNum, err := strconv.Atoi(idSlice[1])
		if err != nil {
			log.Fatal(err)
			return
		}
		idNum++
		return "bug" + "-" + strconv.Itoa(idNum)
	}
	return "bug-0100"
}

//CreateUserID creates a unique user identification number.
func CreateUserID() (OrgID string) {
	return
}

// Encrypt hashes plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}

func session(writer http.ResponseWriter, request *http.Request) (sess Session, err error) {
	cookie, err := request.Cookie("_cookie")
	if err != nil {
		http.Error(writer, "could not retrieve session", http.StatusBadRequest)
	}
	sess = Session{UserID: cookie.Value}
	if ok, _ := sess.Check(); !ok {
		err = errors.New("Invalid session")
		return
	}
	return
}
