package likes

import (
	"database/sql"
	"fmt"
)

type LikesData struct {
	Data *sql.DB
}

func CreateLikesTable(db *sql.DB) *LikesData {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "likes"
	"postID" TEXT NOT NULL
	"username" TEXT NOT NULL,
	"like" TEXT
	
);
`)
	stmt.Exec()
	return &LikesData{
		Data: db,
	}
}


func (likes *LikesData) GetOne(id, user string) LikesFields {
	sliceOfLikeRows := LikesFields{}

	s := fmt.Sprintf("SELECT * FROM likes WHERE postid = '%v' AND username = '%v'", id, user)
	rows, _ := likes.Data.Query(s)
	var postid string
	var author string
	var like string
	if rows.Next() {
		rows.Scan(&postid, &author, &like)
		sliceOfLikeRows = LikesFields{
			PostId:   postid,
			Username: author,
			Like:     like,
		}
	}
	rows.Close()
	return sliceOfLikeRows
}


func (likes *LikesData) Add(postLiked LikesFields) {
	LikedPost := likes.GetOne(postLiked.PostId, postLiked.Username)
	var s string
	if LikedPost.Like == "" {
		s = "INSERT INTO likes (like, postid, username) values (?, ?, ?)"
	} else if postLiked.Like != LikedPost.Like {
		s = "UPDATE likes SET like = ? WHERE postid = ? AND username = ?"
	} else {
		s = "DELETE FROM likes WHERE like = ? AND postid = ? AND username = ?"
	}
	stmt, _ := likes.Data.Prepare(s)
	_, err := stmt.Exec(postLiked.Like, postLiked.PostId, postLiked.Username)
	if err != nil{
		fmt.Println(err)
	}
}

func (likes *LikesData) Get(id, l string) []LikesFields {
	sliceOfLikedRows:= []LikesFields{}
	var s string
	if l == "all" {
		s = fmt.Sprintf("SELECT * FROM likes WHERE username = '%v'", id)

	} else {
		s = fmt.Sprintf("SELECT * FROM likes WHERE postid = '%v' AND like = '%v'", id, l)

	}

	rows, _ := likes.Data.Query(s)
	var postid string
	var author string
	var like string
	for rows.Next() {
		rows.Scan(&postid, &author, &like)
		likedRows := LikesFields{
			PostId:   postid,
			Username: author,
			Like:     like,
		}
		sliceOfLikedRows = append(sliceOfLikedRows, likedRows)
	}
	rows.Close()
	return sliceOfLikedRows
}
