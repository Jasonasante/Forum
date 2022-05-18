package comments

type CommentFields struct {
	CommentId     string
	PostId        string
	Author        string
	Content       string
	CommentAuthor bool
	Likes         int
	Dislikes      int
}
