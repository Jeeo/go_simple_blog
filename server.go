package main

import (
	"net/http"
)

func main() {
	PORT := "3000"
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))

	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/login", loginPage)
	mux.HandleFunc("/signup", signupPage)
	mux.HandleFunc("/register", signup)
	mux.HandleFunc("/authenticate", authenticate)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/userPanel", userPanel)
	mux.HandleFunc("/writePost", createPage)
	mux.HandleFunc("/createPost", createPost)
	mux.HandleFunc("/deletePost/", deletePost)
	mux.HandleFunc("/post/", readPage)

	server := &http.Server{
		Addr:    "0.0.0.0:" + PORT,
		Handler: mux,
	}
	server.ListenAndServe()
}
