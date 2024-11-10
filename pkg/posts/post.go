package posts

import (
	"awesomeProject/pkg/db"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"net/http"
)

type Post struct {
	ID int `db:"id" json:"ID"`

	Title    string         `db:"title" json:"title"`
	Content  string         `db:"content" json:"content"`
	Author   string         `db:"author" json:"author"`
	Visible  bool           `db:"visible" json:"visible"`
	Note     string         `db:"note" json:"note"`
	Tags     pq.StringArray `db:"tags" json:"tags"`
	Category pq.StringArray `db:"category" json:"category"`
	Likes    int
	Views    int

	Username string `db:"username" json:"username"`
}

func InsertPostToDB(db *sqlx.DB, post Post, jwt string) int {
	postID := 0
	post.Username = jwt
	err := db.Get(&postID, `
    INSERT INTO hashed.posts (title, content, author, visible, username, tags, category, note) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
    RETURNING id`,
		post.Title,
		post.Content,
		post.Author,
		post.Visible,
		post.Username,
		post.Tags,
		post.Category,
		post.Note)
	if err != nil {
		return -1
	}
	return postID
}

func CreatePost(post Post, jwt string) int {
	id := InsertPostToDB(db.Conn, post, jwt)
	if id == -1 {
		return -1
	}
	return id
}

func ReadPost(post *Post) (Post, error) {
	err := db.Conn.Get(post, "SELECT * FROM hashed.posts WHERE id=$1", post.ID)
	if err != nil {
		return Post{}, err
	}

	return *post, nil
}

func UpdatePost(post Post) error {
	if post.ID == 0 {
		return fmt.Errorf("ID is not set")
	}
	query := `
    UPDATE hashed.posts
    SET 
        title = COALESCE(NULLIF($1, ''), title),
        content = COALESCE(NULLIF($2, ''), content),
        author = COALESCE(NULLIF($3, ''), author),
        visible = COALESCE($4, visible),
        tags = COALESCE(NULLIF($5::text[], ARRAY[]::text[]), tags),
        category = COALESCE(NULLIF($6::text[], ARRAY[]::text[]), category),
        note = COALESCE(NULLIF($7, ''), note)
    WHERE id = $8;
`

	_, err := db.Conn.Exec(query, post.Title, post.Content, post.Author, post.Visible, post.Tags, post.Category, post.Note, post.ID)
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
