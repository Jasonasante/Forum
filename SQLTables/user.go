package SQLTables

import (
	"database/sql"
	"log"
)

func CreateUserTable(db *sql.DB) {
	createUserSQLTable := "CREATE TABLE IF NOT EXISTS user(userID integer NOT NULL PRIMARY KEY AUTOINCREMENT,username varchar(50) UNIQUE NOT NULL,password varchar(50) NOT NULL,image TEXT)"
	log.Println("Create user table...")
	statement, err := db.Prepare(createUserSQLTable) //prepare SQL statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("user table created")
}

func InsertUser(db *sql.DB, username string, password string, image string) {
	log.Println("Inserting user record ...")
	insertUserSQL := `INSERT INTO user(username, password, image)
	VALUES(?, ?, ?)`
	statement, err := db.Prepare(insertUserSQL) // prepare SQL statement
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(username, password, image)
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
		var password string
		row.Scan(&userID, &username, &password)
		log.Println("User: ", userID, "", username, "", password, "")
	}
}
