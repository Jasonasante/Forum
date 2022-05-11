// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"time"

// 	//"forum/sessions"
// 	"forum/cookies"

// 	_ "github.com/mattn/go-sqlite3"
// 	"golang.org/x/crypto/bcrypt"
// )

// type User struct {
// 	User  string
// 	Email string
// 	Fname string
// 	Lname string
// }

// type ErrorMes struct {
// 	En interface{}
// 	Em string
// }

// var (
// 	database *sql.DB
// 	stmt     *sql.Stmt
// )

// // this receives a password and encrypts it, protect a user's password in the database.
// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	return string(bytes), err
// }

// // this checks whether the inputted string when trying to login matches the encrypted password
// func CheckPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }

// // this determines wither an E-mail exists in the database
// func emailExists(email string) bool {
// 	row := database.QueryRow("SELECT email from user where email= ?", email)
// 	temp := ""
// 	row.Scan(&temp)
// 	if temp != "" {
// 		return true
// 	}
// 	return false
// }

// // this determines wither an username exists in the database
// func usernameExists(username string) bool {
// 	row := database.QueryRow("SELECT username from user where username= ?", username)
// 	temp := ""
// 	row.Scan(&temp)
// 	if temp != "" {
// 		return true
// 	}
// 	return false
// }

// // this sends the inputs in the registration from to the username handleFunc.
// func signUp(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/signup" {
// 		log.Fatal()
// 	}
// 	t, _ := template.ParseFiles("./templates/signup.html")
// 	t.Execute(w, nil)
// }

// // this receives input from the sign up page and inserts the new user information into the database if the username and email does not exist already.
// func avatar(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/avatar" {
// 		log.Fatal()
// 	}
// 	r.ParseForm()

// 	usernameFromSignUp := (r.FormValue("username"))
// 	email := template.HTMLEscapeString(r.FormValue("email"))
// 	password := template.HTMLEscapeString(r.FormValue("psw"))
// 	fname := template.HTMLEscapeString(r.FormValue("fname"))
// 	lname := template.HTMLEscapeString(r.FormValue("lname"))

// 	if emailExists(email) == true && usernameExists(usernameFromSignUp) == false {
// 		en := http.StatusConflict
// 		em := "Uh oh Try again, email already exists!"
// 		t, _ := template.ParseFiles("./templates/errorSignUp.html")
// 		t.Execute(w, ErrorMes{En: en, Em: em})
// 	} else if usernameExists(usernameFromSignUp) == true && emailExists(email) == false {
// 		en := http.StatusConflict
// 		em := "Uh oh Try again, username already exists!"
// 		t, _ := template.ParseFiles("./templates/errorSignUp.html")
// 		t.Execute(w, ErrorMes{En: en, Em: em})
// 	} else if emailExists(email) == true && usernameExists(usernameFromSignUp) == true {
// 		en := http.StatusConflict
// 		em := "Uh oh Try again, username and email already exists!"
// 		t, _ := template.ParseFiles("./templates/errorSignUp.html")
// 		t.Execute(w, ErrorMes{En: en, Em: em})
// 	} else {
// 		stmt, _ = database.Prepare("INSERT INTO user(username,email,password, fname, lname) VALUES(?,?,?,?,?)")
// 		hash, _ := HashPassword(password)
// 		stmt.Exec(usernameFromSignUp, email, hash, fname, lname)
// 		dt := time.Now()
// 		fmt.Print(usernameFromSignUp, " successfully registered ")
// 		fmt.Println("Access granted at", dt.String())

// 		t, err := template.ParseFiles("./templates/avatar.html")
// 		if err != nil {
// 			log.Fatal()
// 		}
// 		t.Execute(w, nil)
// 	}
// }

// // this sends the inputs from the log in form to the homePage handleFunc.
// func logIn(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/login" {
// 		log.Fatal()
// 	}
// 	t, _ := template.ParseFiles("./templates/login.html")
// 	t.Execute(w, nil)
// }

// // this receives the inputs from the log in page and confirms whether the user exists and if the password matches to what is stored in the database determining access to the web page.
// func homePage(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/homepage" {
// 		log.Fatal()
// 	}

// 	r.ParseForm()
// 	usernameFromLogin := template.HTMLEscapeString(r.FormValue("usernameL"))
// 	passwordFromLogin := template.HTMLEscapeString(r.FormValue("pswL"))
// 	var (
// 		usernameFromUserTable string
// 		emailFromUserTable    string
// 		hashFromUserTable     string
// 		fNameFromUserTable    string
// 		lNameFromUserTable    string
// 	)

// 	// this method returns a single row of the information requested within the query that corresponds with the identification key used (i.e username) if it exists
// 	// It then stores the request information in the corresponding variable addresses. Once we check verify that that user exists and the passwords match,we send user to the homepage.
// 	row := database.QueryRow("SELECT * from user WHERE username= ?", usernameFromLogin)
// 	switch err := row.Scan(&usernameFromUserTable, &emailFromUserTable, &hashFromUserTable, &fNameFromUserTable, &lNameFromUserTable); err {
// 	case sql.ErrNoRows:
// 		fmt.Println("No rows were returned!")
// 	case nil:
// 		fmt.Print(usernameFromUserTable + " Info Found. ")
// 	default:
// 		panic(err)
// 	}
// 	if !CheckPasswordHash(passwordFromLogin, hashFromUserTable) {
// 		em := "Username or Password incorrect !! Please try again"
// 		t, _ := template.ParseFiles("./templates/errorLogin.html")
// 		t.Execute(w, ErrorMes{Em: em})
// 	} else {
// 		fmt.Print("Password Matched! Access granted. ")
// 		dt := time.Now()
// 		fmt.Println("Time of Login:", dt.String())
// 		t, _ := template.ParseFiles("./templates/homePage.html")
// 		t.Execute(w, User{usernameFromUserTable, emailFromUserTable, fNameFromUserTable, lNameFromUserTable})
// 	}
// }

// func unregisteredHomePage(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/" {
// 		log.Fatal()
// 	}
// 	cookies.FetchCookies(w, r)
// 	t, _ := template.ParseFiles("./templates/unregisteredHomePage.html")
// 	t.Execute(w, nil)
	
// }

// // this initialises a test sqlite database and creates a table containing user information.
// func TestDB() {
// 	database, _ = sql.Open("sqlite3", "test.db")
// 	stmt, _ = database.Prepare("CREATE TABLE IF NOT EXISTS user (username TEXT, email TEXT, password TEXT, fname TEXT, lname TEXT)")
// 	stmt.Exec()
// }

// // var globalSessions *sessions.Manager

// // // Then, initialize the session manager
// // func CreateGlobalSession() {
// // 	globalSessions,_ = sessions.NewManager("memory", "gosessionid", 3600)
// // }

// func main() {
// 	TestDB()
// 	// CreateGlobalSession()
// 	http.HandleFunc("/", unregisteredHomePage)
// 	http.HandleFunc("/homepage", homePage)
// 	http.HandleFunc("/avatar", avatar)
// 	http.HandleFunc("/login", logIn)
// 	http.HandleFunc("/signup", signUp)
// 	fmt.Println("Starting Server")
// 	fmt.Println("Please open http://localhost:8080/")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// 	fmt.Println("error")
// }
