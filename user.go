package main

import (
	"log"
	"net/http"
	"time"

)

//User is the in-app representation of a organization
type User struct {
	Username   string
	UserID     string
	Email      string
	Password   string
	BugHistory []Bug
	Devs       []Dev
}

//Session models each user-login session
type Session struct {
	ID        int
	Username  string
	UserID    string
	Email     string
	CreatedAt time.Time
}

//AddNewBug allows user to add/assign a new issue to a dev
func (user *User) AddNewBug(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "failed to read request body", http.StatusInternalServerError)
		return
	}
		statement := "insert into bugs (user, title, bugid, description, created, dev) values ($1, $2, $3, $4, $5, $6)"
		_, err = Db.Exec(statement, user.Username, request.PostFormValue("title"), CreateBugID(), request.PostFormValue("des"), time.Now(), request.PostFormValue("dev_name"))
		if err != nil {
			return
		}
		log.Print("Bug created successfully")
}

//Create saves a new user's info into the database
func (user *User) Create() (err error) {
	statement := "insert into users (username, userid, email, password) values($1, $2, $3, $4) returning userid"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return err
	}
	err = stmt.QueryRow(user.Username, user.UserID, user.Email, user.Password).Scan(&user.UserID)
	if err != nil {
		return err
	}
	log.Print("record entered successfully")
	return nil
}

//CreateSession creates a new user login session
func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (username, userid, email, created_at) values ($1, $2, $3, $4) returning id, username, userid, email, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = stmt.QueryRow(user.Username, user.UserID, user.Email, time.Now()).Scan(&session.ID, &session.Username, &session.UserID, session.Email, session.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}
	return session, nil
}

//CheckSession checks and makes sure that a session exists
func (user *User) CheckSession() (session Session, err error) {
	statement := "select id, username, userid, email, created_at from sessions where userid = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	err = stmt.QueryRow(user.UserID).Scan(session.ID, session.Username, session.UserID, session.Email, session.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}
	return session, nil
}

//Check verifies the validity of a session
func (session *Session) Check() (valid bool, err error) {
	statement := "select id, username, userid, email, created_at from sessions where userid = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	err = stmt.QueryRow(session.UserID).Scan(session.ID, session.Username, session.UserID, session.Email, session.CreatedAt)
	if err != nil {
		valid = false
		return
	}
	if session.ID != 0 {
		valid = true
	}
	return
}

//User retrieves a single user from the Database using the session
func (session *Session) User() (user User, err error) {
	err = Db.QueryRow("select username, email from users where userid=$1", session.UserID).Scan(&user.Username, user.Email)
	if err != nil {
		return
	}
	return user, nil
}
