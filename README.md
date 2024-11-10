# Basic REST API Server for Social Media Posts 

This project is a basic REST API server for managing social media posts. The server provides full CRUD (Create, Read, Update, Delete) functionality, allowing users to interact with posts data such as title, content, author, and visibility status.

## Project Status

This project is currently in development. The following features are not fully implemented:

- **Post Filtering**: Options to filter posts based on criteria like tags, categories, and author.
- **Comprehensive Testing**: Unit and integration tests for improving service reliability.
- **Likes and views**: Users can like and comment


## Project Structure

This API is built using Go and is structured to handle HTTP requests for creating, retrieving, updating, and deleting posts in a social media application.

## Data Model

The API interacts with `Post` objects, which have the following structure:

```go
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

```
# `Create` request might look like this:
```
localhost:8080/create
```
```json
{
    "id": 8,
    "title": "A sadasdsdiadasdasdfsdfsdfdsfniasdasdasdmalism",
    "content": "Mini23nal living.",
    "author": "Et4",
    "visible": false,
    "note": "Learnin5g to livef with less.",
    "tags": ["Minim6alism", "Lifestyle", "Simplficity"],
    "category": ["Lifest7yle", "Self-impr8ovefgment"]
}
