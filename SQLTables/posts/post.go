package posts

import (
	"database/sql"
	"fmt"

	"learn.01founders.co/git/gymlad/forum/SQLTables/likes"
)

type PostData struct {
	Data *sql.DB
}

func CreatePostTable(db *sql.DB) *PostData {
	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "posts" (
			"id"	TEXT NOT NULL UNIQUE,
			"author"	TEXT NOT NULL,
			"content"	TEXT NOT NULL,
			"thread"	TEXT,
			PRIMARY KEY("id")
		);
	`)
	stmt.Exec()
	return &PostData{
		Data: db,
	}
}

func (postData *PostData) Add(postFields PostFields) {
	fmt.Println("Post ADD", postFields)
	stmt, _ := postData.Data.Prepare(`INSERT into "posts" 
	(id,author,content,thread) VALUES (?,?,?,?);`)
	stmt.Exec(postFields.Id, postFields.Author, postFields.Content, postFields.Thread)

}

func (posts *PostData) Get(LD *likes.LikesData) []PostFields {
	sliceOfPostTableRows := []PostFields{}
	rows, _ := posts.Data.Query(`SELECT * FROM "posts"`)
	var id string
	var author string
	var content string
	var thread string

	for rows.Next() {
		rows.Scan(&id, &author, &content, &thread)
		// fmt.Println(id, author)
		postTableRows := PostFields{
			Id:       id,
			Author:   author,
			Content:  content,
			Thread:   thread,
			Likes:    len(LD.Get(id, "l")),
			Dislikes: len(LD.Get(id, "d")),
		}
		sliceOfPostTableRows = append(sliceOfPostTableRows, postTableRows)
	}
	rows.Close()
	return sliceOfPostTableRows
}

func (post *PostData) GetMyPosts(likesData *likes.LikesData, str string) ([]PostFields, []PostFields) {
	s := fmt.Sprintf("SELECT * FROM posts WHERE author = '%v'", str)

	sliceOfPostTableRows := []PostFields{}
	sliceOfLikeRows := []PostFields{}
	likesTableRows := likesData.Get(str, "all")
	rows, _ := post.Data.Query(s)
	var id string
	var author string
	var content string
	var thread string
	for rows.Next() {
		rows.Scan(&id, &author, &content, &thread)
		postTableRows := PostFields{
			Id:       id,
			Author:   author,
			Content:  content,
			Thread:   thread,
			Likes:    len(likesData.Get(id, "l")),
			Dislikes: len(likesData.Get(id, "d")),
		}
		sliceOfPostTableRows = append(sliceOfPostTableRows, postTableRows)
	}
	rows.Close()

	for _, v := range likesTableRows {
		s := fmt.Sprintf("SELECT * FROM posts WHERE id = '%v'", v.PostId)

		rows, _ := post.Data.Query(s)
		var id string
		var author string
		var content string
		var thread string
		var postRows PostFields
		if rows.Next() {
			rows.Scan(&id, &author, &content, &thread)
			postRows = PostFields{
				Id:       id,
				Author:   author,
				Content:  content,
				Thread:   thread,
				Likes:    len(likesData.Get(id, "l")),
				Dislikes: len(likesData.Get(id, "d")),
			}
			sliceOfLikeRows = append(sliceOfLikeRows, postRows)
		}
		rows.Close()
	}
	return sliceOfPostTableRows, sliceOfLikeRows
}

func (post *PostData) Filter(likesData *likes.LikesData, str string) []PostFields {
	s := fmt.Sprintf("SELECT * FROM posts WHERE thread LIKE '%v'", "%"+str+"%")

	sliceOfPostRows := []PostFields{}
	rows, _ := post.Data.Query(s)
	var id string
	var author string
	var content string
	var thread string
	for rows.Next() {
		rows.Scan(&id, &author, &content, &thread)
		postRows := PostFields{
			Id:       id,
			Author:   author,
			Content:  content,
			Thread:   thread,
			Likes:    len(likesData.Get(id, "l")),
			Dislikes: len(likesData.Get(id, "d")),
		}
		sliceOfPostRows = append(sliceOfPostRows, postRows)
	}
	rows.Close()
	return sliceOfPostRows
}

func (post *PostData) Delete(id string) {
	stmt, _ := post.Data.Prepare("DELETE FROM posts WHERE id = ?")
	stmt.Exec(id)
}

func (post *PostData) Update(item PostFields, id string) {
	stmt, _ := post.Data.Prepare(`UPDATE "posts" SET "content" = ?, "thread" = ? WHERE "id" = ?`)
	stmt.Exec(item.Content, item.Thread, id)
}
