package main

import (
	"git.learn.01founders.co/gymlad/forum.git/SQLTables"
	"log"
	"database/sql"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

func main(){
	file, err := os.Create("forumDB.db") //create sqlite-compatible database file 
	if err != nil{
		log.Fatal(err.Error())
	}
	file.Close()
	forumDB, _:= sql.Open("sqlite3", "./forumDB.db") //open the created forum-DB file 
	defer forumDB.Close() // defer closing the database
	// create tables into database
	SQLTables.CreateUserTable(forumDB)
	SQLTables.CreateContent(forumDB)
	
	//insert data into tables
	// SQLTables.InsertUser()
	// SQLTables.InsertContent()

	//display table fields 
	// SQLTables.DisplayUser(forumDB)
	// SQLTables.DisplayContent(forumDB)
}