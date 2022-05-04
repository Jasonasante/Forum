package SQLTables

import (
	"database/sql"
	"log"
)

func CreateContent(db *sql.DB) {
	createContentTable := "CREATE TABLE IF NOT EXISTS content(idContent INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,categories TEXT,comments TEXT,imagesURL TEXT,likes INTEGER,dislikes INTEGER)"
	log.Println("Create content table ...")
	statement, err := db.Prepare(createContentTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("content table created")
}

func InsertContent(db *sql.DB, categories string, comments string, imagesURL string, likes int, dislikes int, userID int) {
	log.Println("Inserting content record ...")
	insertUserSQL := `INSERT INTO content(categories, comments, imagesURL, likes, dislikes, userID) 
	VALUES(?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertUserSQL) // prepare SQL statement
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(categories, comments, imagesURL, likes, dislikes, userID)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func DisplayContent(db *sql.DB) {
	row, err := db.Query("SELECT * FROM content ORDER BY categories")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var idContent int
		var categories string
		var comments string
		var imagesURL string
		var likes int
		var dislikes int
		row.Scan(&idContent, &categories, &comments, &imagesURL, &likes, &dislikes)
		log.Println("Content: ", idContent, "", categories, "", comments, "", imagesURL, "", likes, "", dislikes, "")
	}
}
