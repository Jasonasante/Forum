package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

/*type signUp struct {
	Email    string
	Password string
}*/

var (
	database *sql.DB
	stmt     *sql.Stmt
)

func emailExists(email string) bool {
	row := database.QueryRow("SELECT email from user where email= ?", email)
	temp := ""
	row.Scan(&temp)
	if temp != "" {
		return true
	}
	return false
}

func usernameExists(username string) bool {
	row := database.QueryRow("SELECT username from user where username= ?", username)
	temp := ""
	row.Scan(&temp)
	if temp != "" {
		return true
	}
	return false
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		log.Fatal()
	}

	t, _ := template.ParseFiles("./templates/signup.html")
	t.Execute(w, nil)
}

func username(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/username" {
		log.Fatal()
	}
	r.ParseForm()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("psw")
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")

	if emailExists(email) == true && usernameExists(username) == false {
		fmt.Fprint(w, "Uh oh Try again, email already exists!")
	} else if usernameExists(username) == true && emailExists(email) == false {
		fmt.Fprint(w, "Uh oh Try again, username already exists!")
	} else if emailExists(email) == true && usernameExists(username) == true {
		fmt.Fprint(w, "Uh oh Try again, username and email already exists!")
	} else {
		stmt, _ = database.Prepare("INSERT INTO user(username,email,password, fname, lname) VALUES(?,?,?,?,?)")
		stmt.Exec(username, email, password, fname, lname)
	}

	row, _ := database.Query("SELECT * from user")
	for row.Next() {
		var (
			u string
			e string
			p string
			f string
			l string
		)

		row.Scan(&u, &e, &p, &f, &l)
		fmt.Println("username:= " + u + " email:= " + e + " password:= " + p + " fname:= " + f + " lname:= " + l)
	}
	t, err := template.ParseFiles("./templates/avatar.html")
	if err != nil {
		log.Fatal()
	}
	t.Execute(w, nil)
}

func logIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		log.Fatal()
	}
	t, _ := template.ParseFiles("./templates/login.html")
	t.Execute(w, nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/homePage" {
		log.Fatal()
	}
	r.ParseForm()
	email_user := r.FormValue("usernameL")
	password := r.FormValue("pswL")
	var (
		u string
		p string
	)
	if usernameExists(email_user) {
		row, _ := database.Query("SELECT username from user where username=? OR", email_user)
		for row.Next() {
			row.Scan(&u, &p)
			fmt.Println("username:= " + u + " password:= " + p)
		}
	} else {
		fmt.Fprint(w, "username does not exist!")
	}
	if p == password {
		t, _ := template.ParseFiles("./templates/homePage.html")
		t.Execute(w, nil)
	} else {
		fmt.Fprint(w, "Password is incorrect")
	}
}

func main() {
	database, _ = sql.Open("sqlite3", "test.db")
	stmt, _ = database.Prepare("CREATE TABLE IF NOT EXISTS user (username TEXT, email TEXT, password TEXT, fname TEXT, lname TEXT)")
	stmt.Exec()
	http.HandleFunc("/username", username)
	http.HandleFunc("/login", logIn)
	http.HandleFunc("/signup", signUp)
	fmt.Println("hello")
	log.Fatal(http.ListenAndServe(":8081", nil))
	fmt.Println("error")
}
