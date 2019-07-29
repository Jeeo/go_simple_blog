package data

import "log"

type Post struct {
	ID      int
	Author  string
	Title   string
	Content string
}

//Create creates a post
func (post *Post) Create() (err error) {
	statement := "insert into Post (author, title, content) values ($1, $2, $3) returning *"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	err = stmt.QueryRow(post.Author, post.Title, post.Content).Scan(&post.ID, &post.Author, &post.Title, &post.Content)
	if err != nil {
		return
	}
	return
}

//Get Gets a post providing a ID
func (post *Post) Get(id string) (err error) {
	statement := "select Post.id, Post.Title, Post.Content, Users.username from Post, Users where Post.id = $1 and Users.id = Post.Author"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	err = stmt.QueryRow(id).Scan(&post.ID, &post.Title, &post.Content, &post.Author)
	if err != nil {
		return
	}
	return
}

//GetUserPosts Gets all posts written by a user
func GetUserPosts(id string) (posts []Post, err error) {
	statement := "select Post.id, Users.username, Post.Title, Post.Content from Post, Users where Post.Author = $1 AND Users.id = Post.Author"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		log.Println(err)
		return
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content); err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

//GetAll Gets all posts
func GetAll() (posts []Post, err error) {
	rows, err := DB.Query("select Post.id, Users.username, Post.Title, Post.Content  from Post, Users where Users.id = Post.Author")
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content); err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

//Update updates a Post
func (post *Post) Update() (err error) {
	statement := "update Post set author = $1, title = $2, content = $3 where id = $4"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	_, err = stmt.Exec(post.Author, post.Title, post.Title, post.ID)
	if err != nil {
		return
	}
	return
}

// Delete deletes a Post
func (post *Post) Delete() (err error) {
	statement := "delete from Post where id = $1"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	_, err = stmt.Exec(post.ID)
	if err != nil {
		return
	}
	return
}
