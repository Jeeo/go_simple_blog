package main

import (
	"blog/data"
	"log"
	"net/http"
	"strings"
)

func createPage(wri http.ResponseWriter, req *http.Request) {
	auth := checkSession(req)
	if !auth {
		http.Redirect(wri, req, "", 403)
		return
	}

	genPage(wri, struct{ Auth bool }{Auth: auth}, "layout", "navbar", "create_post")
}

func readPage(wri http.ResponseWriter, req *http.Request) {
	auth := checkSession(req)
	id := strings.TrimPrefix(req.URL.Path, "/post/")
	post := &data.Post{}

	if err := post.Get(id); err != nil {
		log.Fatalln(err)
		genPage(wri, struct{ Auth bool }{Auth: auth}, "layout", "navbar", "post")
		return
	}

	genPage(wri, struct {
		Auth bool
		Post data.Post
	}{Auth: auth, Post: *post}, "layout", "navbar", "post")
}

func createPost(wri http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	userIDCookie, err := req.Cookie("USER_ID")
	if err != nil {
		log.Fatalln(err)
		return
	}
	userID := userIDCookie.Value
	post := &data.Post{
		Author:  userID,
		Title:   req.PostFormValue("title"),
		Content: req.PostFormValue("text"),
	}
	err = post.Create()

	if err != nil {
		log.Fatalln("Error on save post: ", err)
		return
	}

	http.Redirect(wri, req, "/", 302)
}

func deletePost(wri http.ResponseWriter, req *http.Request) {
	auth := checkSession(req)
	var user data.User
	post := &data.Post{}
	userSessionCookie, err := req.Cookie("USER_ID")
	if err != nil {
		log.Fatalln(err)
		genPage(wri, struct{ Auth bool }{Auth: auth}, "layout", "navbar", "index")
		return
	}

	postID := strings.TrimPrefix(req.URL.Path, "/deletePost/")
	user, err = data.GetByUserSessionID(userSessionCookie.Value)
	if err != nil {
		log.Println(err)
		genPage(wri, struct{ Auth bool }{Auth: auth}, "layout", "navbar", "user_panel")
		return
	}

	if err := post.Get(postID); err != nil {
		log.Println(err)
		genPage(wri, struct{ Auth bool }{Auth: auth}, "layout", "navbar", "user_panel")
		return
	}
	if post.Author != user.Username {
		log.Println("Not allowed")
		http.Redirect(wri, req, "/", 403)
		return
	}
	if err := post.Delete(); err != nil {
		log.Fatalln(err)
		genPage(wri, struct{ Auth bool }{Auth: auth}, "layout", "navbar", "user_panel")
		return
	}

	genPage(wri, struct {
		Message string
		Auth    bool
	}{Message: "Post deleted", Auth: auth}, "layout", "navbar", "user_panel")

}
