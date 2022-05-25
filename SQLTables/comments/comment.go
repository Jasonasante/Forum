package comments

import (
	"database/sql"
	"fmt"
	"forum/SQLTables/commentsAndLikes"
)

type CommentData struct {
	Data *sql.DB
}

func NewCommentTable(db *sql.DB) *CommentData {
	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "comments" (
			"commentid" TEXT NOT NULL UNIQUE,
			"postid"	TEXT NOT NULL,
			"author"	TEXT NOT NULL,
			"content"	TEXT NOT NULL
		);
	`)
	stmt.Exec()
	return &CommentData{
		Data: db,
	}
}

func (comment *CommentData) Add(commentFields CommentFields) {
	stmt, _ := comment.Data.Prepare(`INSERT INTO "comments" (commentid, postid, author, content) values(?, ?, ?, ?)`)
	_, err := stmt.Exec(commentFields.CommentId, commentFields.PostId, commentFields.Author, commentFields.Content)
	fmt.Println(commentFields)
	fmt.Println(err)
}

func (comment *CommentData) Get(commentsLikesData *commentsAndLikes.CommentsAndLikesData, str string) []CommentFields {
	s := fmt.Sprintf("SELECT * FROM comments WHERE postid = '%v'", str)

	sliceOfCommentRows:= []CommentFields{}
	rows, _ := comment.Data.Query(s)
	var commentid string
	var postid string
	var author string
	var content string
	for rows.Next() {
		rows.Scan(&commentid, &postid, &author, &content)
		commentRows := CommentFields{
			CommentId: commentid,
			PostId:    postid,
			Author:    author,
			Content:   content,
			Likes:         len(commentsLikesData.Get(commentid, "l")),
			Dislikes:         len(commentsLikesData.Get(commentid, "d")),
		}
		sliceOfCommentRows = append(sliceOfCommentRows, commentRows)
	}
	rows.Close()
	return sliceOfCommentRows
}

func (comment *CommentData) Delete(id string) {
	stmt, _ := comment.Data.Prepare(`DELETE FROM "comments" WHERE "commentid" = ?`)
	stmt.Exec(id)
}