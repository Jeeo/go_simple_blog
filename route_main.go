package main

import (
	"blog/data"
	"log"
	"net/http"
)

func index(wri http.ResponseWriter, req *http.Request) {
	auth := checkSession(req)
	allPosts, err := data.GetAll()
	if err != nil {
		log.Fatalln(err)
	}
	genPage(wri, struct {
		Message string
		Auth    bool
		Posts   []data.Post
	}{Message: "", Auth: auth, Posts: allPosts}, "layout", "navbar", "index")
	return
}
