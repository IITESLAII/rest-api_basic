package main

import (
	"awesomeProject/pkg/db"
	jwtToken "awesomeProject/pkg/jwt"
	"awesomeProject/pkg/posts"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func HandlerRegister(dbUser *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		username := queryParams.Get("username")
		password := queryParams.Get("password")
		err := db.InsertToDB(dbUser, db.User{
			Username: username,
			Password: password,
		})
		if err != nil {
			log.Println("Error while inserting user to DB: %v", err)
			err := fmt.Sprintf("Error while inserting user to DB: %v", err)
			w.Write([]byte(err))
		}
	}
}
func HandlerLogin(dbUser *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		username := queryParams.Get("username")
		password := queryParams.Get("password")
		result, b := db.LoginUser(dbUser, db.User{
			Username: username,
			Password: password,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "JWTtoken",
			Value:    *result,
			MaxAge:   4545335,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		if b == false {
			fmt.Fprint(w, "Error while login user")
		} else {
			fmt.Fprint(w, "User logged in")
		}
	}
}

func MiddlewareCheckedExpired(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("JWTtoken")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return

		}
		_, err = jwtToken.ParseJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Error while parsing JWT", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	dbUser := db.Conn
	http.HandleFunc("/register", HandlerRegister(dbUser))
	http.HandleFunc("/login", HandlerLogin(dbUser))
	http.HandleFunc("/create", MiddlewareCheckedExpired(http.HandlerFunc(HandleCreatePost)))
	http.HandleFunc("/read", MiddlewareCheckedExpired(http.HandlerFunc(HandleReadPost)))
	http.HandleFunc("/update", MiddlewareCheckedExpired(http.HandlerFunc(HandleUpdatePost)))
	http.HandleFunc("/delete", MiddlewareCheckedExpired(http.HandlerFunc(HandleDeletePost)))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	post := posts.ParsePost(r)
	err := posts.CreatePost(*post)
	if err != nil {
		log.Println("Error while creating post: %v", err)
	}
}

func HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	post := posts.ParsePost(r)
	err := posts.UpdatePost(*post)
	if err != nil {
		log.Println("Error while updating post: %v", err)
	}
}

func HandleReadPost(w http.ResponseWriter, r *http.Request) {
	post := posts.ParsePost(r)
	s, err := posts.ReadPost(*post)
	if err != nil {
		log.Println("Error while reading post: %v", err)
	}
	fmt.Fprint(w, s)
}

func HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	post := posts.ParsePost(r)
	posts.DeletePost(*post)
}
