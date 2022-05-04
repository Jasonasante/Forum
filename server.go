// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"html/template"
// 	"log"
// 	"net/http"

// 	_ "github.com/mattn/go-sqlite3"
// )

// /*type signUp struct {
// 	Email    string
// 	Password string
// }*/

// var (
// 	database *sql.DB
// 	stmt     *sql.Stmt
// )

// func signUp(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/signup" {
// 		log.Fatal()
// 	}

// 	t, _ := template.ParseFiles("./templates/signup.html")
// 	t.Execute(w, nil)
// }

// func Username(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/username" {
// 		log.Fatal()
// 	}
// 	r.ParseForm()
// 	email := r.FormValue("email")
// 	password := r.FormValue("psw")
// 	stmt, _ = database.Prepare("INSERT INTO user(email,password) VALUES(?,?)")
// 	stmt.Exec(email, password)
// 	row,_:=database.Query("SELECT * from user")
// 	for row.Next(){
// 		var e string
// 		var p string
// 		row.Scan(&e, &p)
// 		fmt.Println("email:= "+ e +" password:= "+ p)
// 	}
// 	t,_:=template.ParseFiles("./template/username.html")
// 	t.Execute(w,nil)
// }

// func main() {
// 	database, _ = sql.Open("sqlite3", "test.db")
// 	stmt, _ = database.Prepare("CREATE TABLE IF NOT EXISTS user (email TEXT, password TEXT)")
// 	stmt.Exec()
// 	http.HandleFunc("/username", Username)
// 	http.HandleFunc("/signup", signUp)
// 	http.ListenAndServe(":8080", nil)
// }
