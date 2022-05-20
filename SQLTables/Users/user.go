package users

import (
	"database/sql"
	"fmt"
)

type UserData struct {
	Data *sql.DB
}

func CreateUserTable(db *sql.DB) *UserData {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "user" (
		"username"	TEXT UNIQUE,
		"email"	TEXT NOT NULL UNIQUE,
		"password"	TEXT NOT NULL,
		"image" TEXT NOT NULL,
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
	INSERT INTO "user" (email, username, password, image) values (?, ?, ?, ?)
	`)
	_, err := stmt.Exec(userFields.Email, userFields.Username, userFields.Password, userFields.Image)
	if err != nil {
		return err
	}
	return nil
}

func (user *UserData) Get() []UserFields {
	sliceOfUserTableRows := []UserFields{}
	rows, _ := user.Data.Query(`
	SELECT * FROM "user"
	`)
	var email string
	var username string
	var password string
	var image string
	for rows.Next() {
		rows.Scan(&email, &username, &password, &image)
		userTableRows := UserFields{
			Email:    email,
			Username: username,
			Password: password,
			Image:    image,
		}
		sliceOfUserTableRows = append(sliceOfUserTableRows, userTableRows)
	}
	rows.Close()
	return sliceOfUserTableRows
}

func (user *UserData) GetUser(str string) UserFields {
	s := fmt.Sprintf("SELECT * FROM user WHERE username = '%v'", str)
	rows, _ := user.Data.Query(s)
	var email string
	var username string
	var password string
	var image string
	var userTableRows UserFields
	if rows.Next() {
		rows.Scan(&email, &username, &password, &image)
		userTableRows = UserFields{
			Email:    email,
			Username: username,
			Password: password,
			Image:    image,
		}
	}
	rows.Close()
	return userTableRows
}

func (user *UserData) UpdateUsername(username string, userFields *UserFields) {
	stmt, _ := user.Data.Prepare(`UPDATE "user" SET "username" = ?,)`)
	_, err := stmt.Exec(username)
	if err != nil {
		fmt.Println(err)
	}
	userFields.Username = username
}

func (user *UserData) UpdateEmail(email string, userFields *UserFields) {
	stmt, _ := user.Data.Prepare(`UPDATE "user" SET "email" = ?,)`)
	_, err := stmt.Exec(email)
	if err != nil {
		fmt.Println(err)
	}
	userFields.Email = email
}

func (user *UserData) UpdatePassword(password string, userFields *UserFields) {
	stmt, _ := user.Data.Prepare(`UPDATE "user" SET "password" = ?,)`)
	_, err := stmt.Exec(password)
	if err != nil {
		fmt.Println(err)
	}
	userFields.Password = password
}

func (user *UserData) DeleteUser(userFields *UserFields) {
	stmt, _ := user.Data.Prepare(`DELETE FROM "user WHERE "username" = ?, password = ?, "email" = ?`)
	_, err := stmt.Exec(userFields.Username, userFields.Password, userFields.Email)
	if err != nil {
		fmt.Println(err)
	}
	userFields = nil
}
