package main

import (
	"blog/data"
	"log"
	"net/http"
)

func userPanel(w http.ResponseWriter, r *http.Request) {
	auth := checkSession(r)
	ProvidedData := struct {
		Message string
		Auth    bool
		Posts   []data.Post
	}{
		Message: "", Auth: auth, Posts: []data.Post{}}
	cookieUserID, err := r.Cookie(COOKIE_USER_ID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", 403)
		return
	}
	posts, err := data.GetUserPosts(cookieUserID.Value)
	if err != nil || len(posts) < 1 {
		log.Println(err)
		ProvidedData.Message = "You don't have written a post yet :("
		genPage(w, ProvidedData, "layout", "navbar", "user_panel")
		return
	}
	ProvidedData.Message = "Your Posts"
	ProvidedData.Posts = posts
	genPage(w, ProvidedData, "layout", "navbar", "user_panel")
}
