package SQLTables

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func CreateUserTable(db *sql.DB) {
	createUserSQLTable :=  `CREATE TABlE IF NOT EXISTS user(id INTEGER PRIMARY KEY, username TEXT, email TEXT UNIQUE, password TEXT, firstname TEXT, lastname TEXT)` // SQL Statement for Create Table

	log.Println("Create user table...")
	statement, err := db.Prepare(createUserSQLTable) //prepare SQL statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("user table created")
}

func InsertUser(db *sql.DB, username string,email string, password string,firstname string,lastname string) {
	log.Println("Inserting user record ...")
	insertUserSQL := `INSERT INTO user(username, email, password, firstname, lastname) 
	VALUES(?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertUserSQL) // prepare SQL statement
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(username, email, password, firstname, lastname)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func DisplayUser(db *sql.DB) {
	row, err := db.Query("SELECT * FROM user ORDER BY username")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var userID int
		var username string
		var email string
		var password string
		var firstname string
		var lastname string
		row.Scan(&userID, &username,&email, &password,&firstname,&lastname)
		log.Println("User: ", userID, "", username, "",email,"", password, "",firstname,"",lastname,"")
	}
}
