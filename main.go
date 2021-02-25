package main

import (
	"net/http"

	"github.com/gorilla/mux"

)

func main() {
	router := mux.NewRouter()
	server := http.Server{
		Addr:    "127.0.0.1:9090",
		Handler: router,
	}
	server.ListenAndServe()
	router.HandleFunc("/", home)
	router.HandleFunc("/newbug", CreateBug)

}

func home(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Welcome"))
}

//UserByEmail retrieves a single user from the database
func UserByEmail(email string) (user User, err error) {
	err = Db.QueryRow("SELECT username, userid, password, FROM users WHERE email = $1", email).Scan(&user.Username, &user.UserID, &user.Password)
	if err != nil {
		return
	}
	return user, nil
}

//CreateBug handles the creation/listing of a new issue/bug
func CreateBug(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err!=nil{
		http.Error(writer, "Invalid session", http.StatusBadRequest)
		return
	}
	user, err := sess.User()
	if err!=nil{
		http.Error(writer, "User Not Found", http.StatusBadRequest)
		return
	}
	user.AddNewBug(writer, request)

}
