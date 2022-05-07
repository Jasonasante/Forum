package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	User string
}

type ErrorMes struct {
	En interface{}
	Em string
}

var (
	database *sql.DB
	stmt     *sql.Stmt
)

//this receives a password and encrypts it, protect a user's password in the database.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

//this checks whether the inputted string when trying to login matches the encrypted password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//this determines wither an E-mail exists in the database
func emailExists(email string) bool {
	row := database.QueryRow("SELECT email from user where email= ?", email)
	temp := ""
	row.Scan(&temp)
	if temp != "" {
		return true
	}
	return false
}

//this determines wither an username exists in the database
func usernameExists(username string) bool {
	row := database.QueryRow("SELECT username from user where username= ?", username)
	temp := ""
	row.Scan(&temp)
	if temp != "" {
		return true
	}
	return false
}

//this sends the inputs in the registration from to the username handleFunc.
func signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		log.Fatal()
	}

	t, _ := template.ParseFiles("./templates/signup.html")
	t.Execute(w, nil)
}

//this receives input from the sign up page and inserts the new user information into the database if the username and email does not exist already.
func avatar(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/avatar" {
		log.Fatal()
	}
	r.ParseForm()

	username := template.HTMLEscapeString(r.FormValue("username"))
	email := template.HTMLEscapeString(r.FormValue("email"))
	password := template.HTMLEscapeString(r.FormValue("psw"))
	fname := template.HTMLEscapeString(r.FormValue("fname"))
	lname := template.HTMLEscapeString(r.FormValue("lname"))

	if emailExists(email) == true && usernameExists(username) == false {
		en := http.StatusConflict
		em := "Uh oh Try again, email already exists!"
		t, _ := template.ParseFiles("./templates/errorSignUp.html")
		t.Execute(w, ErrorMes{En: en, Em: em})
	} else if usernameExists(username) == true && emailExists(email) == false {
		en := http.StatusConflict
		em := "Uh oh Try again, username already exists!"
		t, _ := template.ParseFiles("./templates/errorSignUp.html")
		t.Execute(w, ErrorMes{En: en, Em: em})
	} else if emailExists(email) == true && usernameExists(username) == true {
		en := http.StatusConflict
		em := "Uh oh Try again, username and email already exists!"
		t, _ := template.ParseFiles("./templates/errorSignUp.html")
		t.Execute(w, ErrorMes{En: en, Em: em})
	} else {
		stmt, _ = database.Prepare("INSERT INTO user(username,email,password, fname, lname) VALUES(?,?,?,?,?)")
		hash, _ := HashPassword(password)
		stmt.Exec(username, email, hash, fname, lname)
		dt:=time.Now()
		fmt.Print(username , "successfully registered ")
		fmt.Println("Access granted at", dt.String())
		t, err := template.ParseFiles("./templates/avatar.html")
		if err != nil {
			log.Fatal()
		}
		t.Execute(w, nil)
	}

	// row, _ := database.Query("SELECT * from user")
	// for row.Next() {
	// 	var (
	// 		u string
	// 		e string
	// 		p string
	// 		f string
	// 		l string
	// 	)

	// 	row.Scan(&u, &e, &p, &f, &l)
	// 	fmt.Println("username:= " + u + " email:= " + e + " password:= " + p + " fname:= " + f + " lname:= " + l)
	// }
}

//this sends the inputs from the log in form to the homePage handleFunc.
func logIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		log.Fatal()
	}
	t, _ := template.ParseFiles("./templates/login.html")
	t.Execute(w, nil)
}

//this receives the inputs from the log in page and confirms whether the user exists and if the password matches to what is stored in the database determining access to the web page.
func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/homepage" {
		log.Fatal()
	}
	r.ParseForm()
	username := template.HTMLEscapeString(r.FormValue("usernameL"))
	password := template.HTMLEscapeString(r.FormValue("pswL"))
	var (
		u    string
		hash string
	)
	//this method returns a single row of the information requested within the query that corresponds with the identification key used (i.e username) if it exists
	//It then stores the request information in the corresponding variable addresses. Once we check verify that that user exists and the passwords match,we send user to the homepage.
	row := database.QueryRow("SELECT username, password from user WHERE username= ?", username)
	switch err := row.Scan(&u, &hash); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Print(u + " Info Found. ")
	default:
		panic(err)
	}
	if !CheckPasswordHash(password, hash) {
		em := "Username or Password incorrect !! Please try again"
		t, _ := template.ParseFiles("./templates/errorLogin.html")
		t.Execute(w, ErrorMes{Em: em})
	} else {
		fmt.Print("Password Matched! Access granted. ")
		dt := time.Now()
		fmt.Println("Time of Login:", dt.String())
		t, _ := template.ParseFiles("./templates/homePage.html")
		t.Execute(w, User{User: u})
	}
}

//this initialises a test sqlite database and creates a table containing user information.
func TestDB() {
	database, _ = sql.Open("sqlite3", "test.db")
	stmt, _ = database.Prepare("CREATE TABLE IF NOT EXISTS user (username TEXT, email TEXT, password TEXT, fname TEXT, lname TEXT)")
	stmt.Exec()
}
func main() {
	TestDB()
	http.HandleFunc("/homepage", homePage)
	http.HandleFunc("/avatar", avatar)
	http.HandleFunc("/login", logIn)
	http.HandleFunc("/signup", signUp)
	fmt.Println("Starting Server")
	fmt.Println("Please open http://localhost:8080/signup")
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("error")
}
