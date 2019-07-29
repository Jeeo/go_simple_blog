package main

import (
	"blog/data"
	"fmt"
	"log"
	"net/http"
	"time"
)

// CookieName define a Cookie session name
const (
	COOKIE_SESSION_UUID = "SESSION_UUID"
	COOKIE_USER_ID      = "USER_ID"
)

func loginPage(w http.ResponseWriter, r *http.Request) {
	auth := false
	genPage(w, struct {
		Message string
		Auth    bool
	}{Message: "", Auth: auth}, "layout", "navbar", "login")
}

func signupPage(w http.ResponseWriter, r *http.Request) {
	auth := false
	genPage(w, struct {
		Message string
		Auth    bool
	}{Message: "", Auth: auth}, "layout", "navbar", "register")
}

func signup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := &data.User{
		Username: r.PostFormValue("username"),
		Password: data.Encrypt(r.PostFormValue("password")),
	}
	if err := user.Create(); err != nil {
		log.Println(err.Error())
		genPage(w, struct {
			Message string
			Auth    bool
		}{
			Message: "Error on create your account",
			Auth:    false,
		}, "layout", "navbar", "index")
		return
	}

	welcome := fmt.Sprintf("Welcome %s", user.Username)
	genPage(w, struct {
		Message string
		Auth    bool
	}{
		Message: welcome,
		Auth:    false,
	}, "layout", "navbar", "index")
}

func authenticate(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.PostFormValue("username")
	password := data.Encrypt(req.PostFormValue("password"))
	if user, err := data.Get(username); err != nil {
		log.Println(err)
		genPage(w, struct {
			Message string
			Auth    bool
		}{Message: "USER NOT FOUND", Auth: false}, "layout", "navbar", "index")
		return
	} else if user.Password == password {
		session, err := data.CreateSession(user.ID)
		if err != nil {
			log.Fatal(err, " Error on create session")
			return
		}
		cookieSessionUUID, cookieUserID := handleCookies(session.Uuid, session.UserId, false)
		http.SetCookie(w, cookieSessionUUID)
		http.SetCookie(w, cookieUserID)
		http.Redirect(w, req, "/", 302)
	} else {
		genPage(w, struct {
			Message string
			Auth    bool
		}{Message: "INCORRECT PASSWORD", Auth: false}, "layout", "navbar", "index")
	}
}

func logout(w http.ResponseWriter, req *http.Request) {
	cookieUserID, err := req.Cookie(COOKIE_USER_ID)
	if err != nil {
		log.Fatal(err, " ERROR on obtain session cookie - userid")
		return
	}
	userID := cookieUserID.Value
	cookieSessionUUID, err := req.Cookie(COOKIE_SESSION_UUID)
	if err != nil {
		log.Fatal(err, "ERROR on obtain session cookie - sessionuuid")
		return
	}
	if session, err := data.GetSessionByUserID(userID); err != nil {
		log.Fatal(err, " ERROR on get session to logout")
		return
	} else {
		if err := session.DeleteSessionByUUID(); err != nil {
			log.Fatal(err, " ERROR on delete session in database")
			return
		}

		cookieSessionUUID, cookieUserID = handleCookies("", 0, true) // delete cookies
		http.SetCookie(w, cookieSessionUUID)
		http.SetCookie(w, cookieUserID)
	}
	http.Redirect(w, req, "/", 302)
}

func handleCookies(sessionUUID string, userID int, delete bool) (cookieSessionUUID *http.Cookie, cookieUserID *http.Cookie) {
	minutes := 10
	maxAge := minutes * 60

	cookieSessionUUID = &http.Cookie{
		Name:     COOKIE_SESSION_UUID,
		Value:    sessionUUID,
		HttpOnly: true,
		MaxAge:   maxAge,
	}

	cookieUserID = &http.Cookie{
		Name:     COOKIE_USER_ID,
		Value:    fmt.Sprintf("%d", userID),
		HttpOnly: true,
		MaxAge:   maxAge,
	}

	if delete {
		cookieSessionUUID.Expires = time.Unix(0, 0)
		cookieSessionUUID.MaxAge = 0
		cookieUserID.Expires = time.Unix(0, 0)
		cookieUserID.MaxAge = 0
	}

	return
}

func checkSession(req *http.Request) (auth bool) {
	session, err := req.Cookie("SESSION_UUID")
	if err != nil {
		log.Println(err)
		return
	}

	cookieUUID := session.Value
	dataSession := data.Session{Uuid: cookieUUID}
	if valid, err := dataSession.Check(); err != nil {
		log.Fatalln(err)
	} else if valid {
		return true
	}
	return
}
