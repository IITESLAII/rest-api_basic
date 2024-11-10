# Basic REST API Server for Social Media Posts 

This project is a basic REST API server for managing social media posts. The server provides full CRUD (Create, Read, Update, Delete) functionality, allowing users to interact with posts data such as title, content, author, and visibility status.

## Project Status

This project is currently in development. The following features are not fully implemented:

- **Post Editing Permissions**: Only the author of a post should be able to edit it.
- **Tags and Categories**: Support for tags and categories to organize posts.
- **Notes for Posts**: Ability to add internal notes to posts (for administrative use).
- **Post Filtering**: Options to filter posts based on criteria like tags, categories, and author.
- **Comprehensive Testing**: Unit and integration tests for improving service reliability.


## Project Structure

This API is built using Go and is structured to handle HTTP requests for creating, retrieving, updating, and deleting posts in a social media application.

## Data Model

The API interacts with `Post` objects, which have the following structure:

```go
type Post struct {
    ID      int    `db:"id" json:"ID"`
    Title   string `db:"title" json:"title"`
    Content string `db:"content" json:"content"`
    Author  string `db:"author" json:"author"`
    Visible bool   `db:"visible" json:"visible"`
    User    string `db:"user" json:"user"`
}
