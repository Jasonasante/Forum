package Users

import (
	"database/sql"
	"fmt"
)

type UserData struct {
	Data *sql.DB
}

func CreateUserTable(db *sql.DB) *UserData {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "users" (
		"email"	TEXT UNIQUE,
		"username"	TEXT NOT NULL UNIQUE,
		"password"	TEXT NOT NULL,
		"picture" TEXT NOT NULL,
		PRIMARY KEY("username")
	);
`)
	stmt.Exec()

	return &UserData{
		Data: db,
	}
}

func (user *UserData) Add(userFields UserFields) error {
	stmt, _ := user.Data.Prepare(`
	INSERT INTO "users" (email, username, password, picture) values (?, ?, ?, ?)
	`)
	_, err := stmt.Exec(userFields.Email, userFields.Username, userFields.Password, userFields.Picture)
	if err != nil {
		return err
	}
	return nil
}

func (user *UserData) Get() []UserFields {
	sliceOfUserTableRows := []UserFields{}
	rows, _ := user.Data.Query(`
	SELECT * FROM "users"
	`)
	var email string
	var username string
	var password string
	var picture string 
	for rows.Next() {
		rows.Scan(&email, &username, &password, &picture)
		userTableRows := UserFields{
			Email:    email,
			Username: username,
			Password: password,
			Picture: picture,
		}
		sliceOfUserTableRows = append(sliceOfUserTableRows, userTableRows)
	}
	rows.Close()
	return sliceOfUserTableRows
}

func (user *UserData) GetUser(str string) UserFields {
	s := fmt.Sprintf("SELECT * FROM users WHERE username = '%v'", str)
	rows, _ := user.Data.Query(s)
	var email string
	var username string
	var password string
	var picture string 
	var userTableRows UserFields
	if rows.Next() {
		rows.Scan(&email, &username, &password)
		userTableRows = UserFields{
			Email:    email,
			Username: username,
			Password: password,
			Picture: picture,
		}
	}
	rows.Close()
	return userTableRows
}

func (user *UserData) UpdateUsername(username string, userFields *UserFields) {
	stmt, _ := user.Data.Prepare(`UPDATE "users" SET "username" = ?,)`)
	_, err := stmt.Exec(username)
	if err != nil {
		fmt.Println(err)
	}
	userFields.Username = username
}

func (user *UserData) UpdateEmail(email string, userFields *UserFields) {
	stmt, _ := user.Data.Prepare(`UPDATE "users" SET "email" = ?,)`)
	_, err := stmt.Exec(email)
	if err != nil {
		fmt.Println(err)
	}
	userFields.Email = email
}

func (user *UserData) UpdatePassword(password string, userFields *UserFields) {
	stmt, _ := user.Data.Prepare(`UPDATE "users" SET "password" = ?,)`)
	_, err := stmt.Exec(password)
	if err != nil {
		fmt.Println(err)
	}
	userFields.Password = password
}

func (user *UserData) DeleteUser(userFields *UserFields) {
	stmt, _ := user.Data.Prepare(`DELETE FROM "users WHERE "username" = ?, password = ?, "email" = ?`)
	_, err := stmt.Exec(userFields.Username, userFields.Password, userFields.Email)
	if err != nil {
		fmt.Println(err)
	}
	userFields = nil
}
