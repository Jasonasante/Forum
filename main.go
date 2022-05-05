package main

import (
	"forum/SQLTables"
	// "forum/SQLcontent"
	// "log"
	"database/sql"
	// "os"
	_ "github.com/mattn/go-sqlite3"
)

func main(){
	// file, err := os.Create("forum-DB.db") //create sqlite-compatible database file 
	// if err != nil{
	// 	log.Fatal(err.Error())
	// }
	// file.Close()
	forumDB, _:= sql.Open("sqlite3", "./forum-DB.db") //open the created forum-DB file 
	defer forumDB.Close() // defer closing the database

	// create tables into database
	// SQLUser.CreateUserTable(forumDB) 
	SQLTables.CreateContent(forumDB)
	//insert data into tables
	// SQLUser.InsertUser(forumDB,"miguel","124")
	// insertContent()
	//display table fields 
	SQLTables.DisplayUser(forumDB)
	SQLTables.DisplayContent(forumDB)
}



// stuser, _ := db.Prepare(`CREATE TABlE IF NOT EXISTS user(id INTEGER PRIMARY KEY, username TEXT, email TEXT UNIQUE, password TEXT, first_name TEXT, last_name TEXT)`)
// 	stpost, _ := db.Prepare(`CREATE TABlE IF NOT EXISTS post(id INTEGER PRIMARY KEY, user_id INTEGER, categ_id INTEGER, title TEXT, post TEXT,  FOREIGN KEY(user_id) REFERENCES user(id), FOREIGN KEY(categ_id) REFERENCES categ(id))`)

	