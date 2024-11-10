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
		u, err := db.ParsingLogReg(r)
		if u.Username == "" || u.Password == "" {
			w.Write([]byte("empty username or password"))
			return
		}
		if err != nil {
			log.Fatalf("Error while parsing user: %v", err)
		}
		err = db.InsertToDB(dbUser, u)
		if err != nil {
			log.Println("Error while inserting user to DB: %v", err)
			err := fmt.Sprintf("Error while inserting user to DB: %v", err)
			w.Write([]byte(err))
		}
	}
}
func HandlerLogin(dbUser *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := db.ParsingLogReg(r)
		if err != nil {
			log.Fatalf("Error while parsing user: %v", err)
		}
		if u.Username == "" || u.Password == "" {
			w.Write([]byte("empty username or password"))
			return
		}
		result, b := db.LoginUser(dbUser, u)
		if !b {
			w.Write([]byte("wrong username or password"))
			return
		}
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
	cookie, err := r.Cookie("JWTtoken")
	if err != nil {
		return
	}
	jwt, err := jwtToken.ParseJWT(cookie.Value)
	if err != nil {
		return
	}
	id := posts.CreatePost(*post, jwt)
	if id == -1 {
		log.Println("Error while creating post")
	}
	_, err = fmt.Fprint(w, id)
	if err != nil {
		return
	}
}

func HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	post := posts.ParsePost(r)
	postS := *post
	postForUpdate, err := posts.ReadPost(&postS)
	if err != nil {
		log.Println("Error while reading post: %v", err)
	}
	if !CheckUserValid(r, &postForUpdate) {
		http.Error(w, "no access", http.StatusForbidden)
		return
	}
	err = posts.UpdatePost(*post)
	if err != nil {
		fmt.Fprintf(w, "Error while updating post: %v", err)
	} else {
		fmt.Fprint(w, http.StatusOK)
	}
}

func HandleReadPost(w http.ResponseWriter, r *http.Request) {
	post := posts.ParsePost(r)
	s, err := posts.ReadPost(post)
	if err != nil {
		log.Println("Error while reading post: %v", err)
	}
	if post.Visible == false {
		if !CheckUserValid(r, post) {
			http.Error(w, "no access", http.StatusForbidden)
			return
		}
	}
	fmt.Fprint(w, s)
}

func HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	post := posts.ParsePost(r)
	postS := *post
	postForDelete, err := posts.ReadPost(&postS)
	if err != nil {
		fmt.Errorf("error while reading the post")
	}
	if !CheckUserValid(r, &postForDelete) {
		http.Error(w, "no access", http.StatusForbidden)
		return
	}
	posts.DeletePost(*post)
	fmt.Fprint(w, http.StatusOK)
}

func CheckUserValid(r *http.Request, post *posts.Post) bool {
	cookie, err := r.Cookie("JWTtoken")
	if err != nil {
		fmt.Errorf("error while getting cookie: %v", err)
		return false
	}
	if jwtUserName, _ := jwtToken.ParseJWT(cookie.Value); jwtUserName != post.Username {
		fmt.Errorf("error user is not the author of post")
		return false
	}
	return true
}
