package main

import "net/http"

func signupAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "signup failed", http.StatusBadRequest)
		return
	}
	user:= &User{
		Username:request.PostFormValue("username"),
		UserID:CreateUserID(),
		Email: request.PostFormValue("email"),
		Password:Encrypt(request.PostFormValue("password")),
	}
	err = user.Create()
	if err != nil {
		http.Error(writer, "signup failed", http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/login", 302)
}

func authenticate(writer http.ResponseWriter, request *http.Request) {
	err:= request.ParseForm()
	if err!=nil{
		http.Error(writer, "failed to sign user in", http.StatusBadRequest)
		return
	}
	user, err:= UserByEmail(request.PostFormValue("email"))
	if err!=nil{
	http.Error(writer, "User Not Found", http.StatusBadRequest)
	return
	}
	if user.Password == Encrypt(request.PostFormValue("password")){
		session, err:= user.CreateSession()
		if err!=nil{
			http.Error(writer, "Cannot Create Session", http.StatusInternalServerError)
			return
		}
		cookie:= http.Cookie{
			Name: "_cookie",
			Value: session.UserID,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)
		http.Redirect(writer, request, "/home", 302)
	}else {
 		http.Redirect(writer, request, "/login", 302)
	}
}

