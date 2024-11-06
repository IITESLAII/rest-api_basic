package posts

import (
	"awesomeProject/pkg/db"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type Post struct {
	ID int `db:"id" json:"ID"`

	Title   string `db:"title" json:"title"`
	Content string `db:"content" json:"content"`
	Author  string `db:"author" json:"author"`
	Visible bool   `db:"visible" json:"visible"`

	User string
}

func InsertPostToDB(db *sqlx.DB, post Post) error {
	_, err := db.Exec("INSERT INTO hashed.posts (title, content, author, visible, user) VALUES ($1, $2, $3, $4, $5)", post.Title, post.Content, post.Author, post.Visible, post.User)
	if err != nil {
		return err
	}
	return nil
}

func CreatePost(post Post) error {
	err := InsertPostToDB(db.Conn, post)
	if err != nil {
		return err
	}
	return nil
}

func ReadPost(post Post) (Post, error) {
	id := post.ID
	err := db.Conn.Get(&post, "SELECT * FROM hashed.posts WHERE id=$1", id)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func UpdatePost(post Post) error {
	if post.ID == 0 {
		return fmt.Errorf("ID is not set")
	}
	if post.Title != "" {
		_, err := db.Conn.Exec("UPDATE hashed.posts SET title=$1 WHERE id=$2", post.Title, post.ID)
		if err != nil {
			return err
		}
	}
	if post.Content != "" {
		_, err := db.Conn.Exec("UPDATE hashed.posts SET content=$1 WHERE id=$2", post.Content, post.ID)
		if err != nil {
			return err
		}
	}
	if post.Author != "" {
		_, err := db.Conn.Exec("UPDATE hashed.posts SET author=$1 WHERE id=$2", post.Author, post.ID)
		if err != nil {
			return err
		}
	}

	_, err := db.Conn.Exec("UPDATE hashed.posts SET visible=$1 WHERE id=$2", post.Visible, post.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeletePost(post Post) {
	_, err := db.Conn.Exec("DELETE FROM hashed.posts WHERE id=$1", post.ID)
	if err != nil {
		fmt.Printf("Error while deleting post: %v", err)
	}
}

func ParsePost(r *http.Request) *Post {
	post := Post{}
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		log.Fatalf("Error while parsing post: %v", err)
		return nil
	}
	return &post
}
