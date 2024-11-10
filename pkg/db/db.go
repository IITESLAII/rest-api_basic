package db

import (
	"awesomeProject/hash"
	jwtToken "awesomeProject/pkg/jwt"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type User struct {
	Username string `db:"username"`
	Password string `db:"password_hashed"`
	UUID     string `db:"uuid"`
}

func ConnectToDB() *sqlx.DB {
	dbUser, err := sqlx.Connect("postgres", "user=postgres password=123 dbname=awesomeproject sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	return dbUser
}

var Conn = ConnectToDB()

func CheckUserExist(db *sqlx.DB, username string) bool {
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM hashed.passusers WHERE username=$1)", username)
	if err != nil {
		fmt.Printf("Error while check user exists from DB: %v", err)
	}
	return exists
}

func InsertToDB(db *sqlx.DB, user User) error {
	if CheckUserExist(db, user.Username) {
		return fmt.Errorf("User already exists")
	}
	hashPassword := hash.PasswordHash(user.Password)
	u := uuid.New()
	_, err := db.Exec("INSERT INTO hashed.passusers (username, password_hashed, uuid) VALUES ($1, $2, $3)", user.Username, hashPassword, u)
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(db *sqlx.DB, user User) (*string, bool) {
	var person User
	err := db.Get(&person, "SELECT username, password_hashed, uuid FROM hashed.passusers WHERE username=$1", user.Username)
	if err != nil {
		fmt.Printf("Error while getting user from DB: %v", err)
	}
	if hash.CheckPasswordHash(user.Password, []byte(person.Password)) {
		token, err := jwtToken.GenerateJWT(person.Username, person.UUID)
		if err != nil {
			fmt.Printf("Error while generating JWT: %v", err)
		}
		return &token, true
	}

	return nil, false
}

func ParsingLogReg(r *http.Request) (User, error) {
	u := User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Fatalf("Error while parsing user: %v", err)
		return User{}, nil
	}
	return u, nil
}
