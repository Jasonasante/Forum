package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	users "forum/SQLTables/Users"
	"forum/SQLTables/comments"
	"forum/SQLTables/commentsAndLikes"
	"forum/SQLTables/likes"
	"forum/SQLTables/posts"
	"forum/cookies"
	"forum/sessions"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	UserTable                  *users.UserData
	CommentTable               *comments.CommentData
	LikesDislikesTable         *likes.LikesData
	LikesDislikesCommentsTable *commentsAndLikes.CommentsAndLikesData
	PostsTable                 *posts.PostData
)

type ErrorMes struct {
	En interface{}
	Em string
}

// this receives a password and encrypts it, protect a user's password in the database.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// this checks whether the inputted string when trying to login matches the encrypted password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// this determines wither an E-mail exists in the database
func emailExists(email string) bool {
	row := UserTable.Data.QueryRow("SELECT email from user where email= ?", email)
	temp := ""
	row.Scan(&temp)
	if temp != "" {
		return true
	}
	return false
}

// this determines wither an username exists in the database
func usernameExists(username string) bool {
	row := UserTable.Data.QueryRow("SELECT username from user where username= ?", username)
	temp := ""
	row.Scan(&temp)
	if temp != "" {
		return true
	}
	return false
}

// this sends the inputs in the registration from to the username handleFunc.
func signUp(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if r.URL.Path != "/signup" {
		log.Fatal()
	}

	if AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t, _ := template.ParseFiles("./templates/signup.html")
	t.Execute(w, nil)
}

// this receives input from the sign up page and inserts the new user information into the database if the username and email does not exist already.
func avatar(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if r.URL.Path != "/avatar" {
		log.Fatal()
	}

	r.ParseForm()

	usernameFromSignUp := (r.FormValue("username"))
	email := template.HTMLEscapeString(r.FormValue("email"))
	password := template.HTMLEscapeString(r.FormValue("psw"))

	if emailExists(email) && !usernameExists(usernameFromSignUp) {
		en := http.StatusConflict
		w.WriteHeader(en)
		em := "Uh oh Try again, email already exists!"
		t, _ := template.ParseFiles("./templates/errorSignUp.html")
		t.Execute(w, ErrorMes{En: en, Em: em})
	} else if usernameExists(usernameFromSignUp) && !emailExists(email) {
		en := http.StatusConflict
		w.WriteHeader(en)
		em := "Uh oh Try again, username already exists!"
		t, _ := template.ParseFiles("./templates/errorSignUp.html")
		t.Execute(w, ErrorMes{En: en, Em: em})
	} else if emailExists(email) && usernameExists(usernameFromSignUp) {
		en := http.StatusConflict
		w.WriteHeader(en)
		em := "Uh oh Try again, username and email already exists!"
		t, _ := template.ParseFiles("./templates/errorSignUp.html")
		t.Execute(w, ErrorMes{En: en, Em: em})
	} else {
		hash, _ := HashPassword(password)
		UserTable.Add(users.UserFields{Email: email, Username: usernameFromSignUp, Password: hash})
		dt := time.Now()
		fmt.Print(usernameFromSignUp, " successfully registered ")
		fmt.Println("Access granted at", dt.String())
		cookie := &http.Cookie{
			Name:  "Maryland",
			Value: "1",
		}
		http.SetCookie(w, cookie)
		t, err := template.ParseFiles("./templates/avatar.html")
		if err != nil {
			log.Fatal()
		}
		t.Execute(w, nil)
		s.IsAuthorized = true
		s.Username = usernameFromSignUp
		s.Expiry = time.Now().Add(120 * time.Second)

	}
}

func uploadFile(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	c, err := r.Cookie("sessionId")
	if err != nil {
		fmt.Println("11")
	}
	session := sessions.SessionMap.Get(c.Value)
	DpFile, err := os.Create("dp-images/" + session.Username + "-dp.png")
	DpName := "../dp-images/" + session.Username + "-dp.png"
	UserTable.Data.Exec("UPDATE user SET image = ? WHERE username = ?", DpName, session.Username)
	if err != nil {
		fmt.Println(err)
	}
	defer DpFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	DpFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	// fmt.Fprintf(w, "Successfully Uploaded File\n")
	http.Redirect(w, r, "/", http.StatusFound)
}

// this sends the inputs from the log in form to the homePage handleFunc.
func logIn(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if r.URL.Path != "/login" {
		log.Fatal()
	}
	if AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t, _ := template.ParseFiles("./templates/login.html")
	t.Execute(w, nil)
}

// this receives the inputs from the log in page and confirms whether the user exists and if the password matches to what is stored in the database determining access to the web page.
func AuthoriseLogin(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	r.ParseForm()
	usernameFromLogin := template.HTMLEscapeString(r.FormValue("usernameL"))
	passwordFromLogin := template.HTMLEscapeString(r.FormValue("pswL"))
	var (
		usernameFromUserTable string
		emailFromUserTable    string
		hashFromUserTable     string
		iFromUserTable        string
	)

	// this method returns a single row of the information requested within the query that corresponds with the identification key used (i.e username) if it exists
	// It then stores the request information in the corresponding variable addresses. Once we check verify that that user exists and the passwords match,we send user to the homepage.
	row := UserTable.Data.QueryRow("SELECT * from user WHERE username= ?", usernameFromLogin)
	switch err := row.Scan(&usernameFromUserTable, &emailFromUserTable, &hashFromUserTable, &iFromUserTable); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Print(usernameFromUserTable + " Info Found. ")
	default:
		panic(err)
	}
	if !CheckPasswordHash(passwordFromLogin, hashFromUserTable) {
		en := http.StatusBadRequest
		w.WriteHeader(en)
		em := "Username or Password incorrect !! Please try again"
		t, _ := template.ParseFiles("./templates/errorLogin.html")
		t.Execute(w, ErrorMes{En: en, Em: em})
	} else {
		fmt.Print("Password Matched! Access granted. ")
		dt := time.Now()
		fmt.Println("Time of Login:", dt.String())
		for k, v := range sessions.SessionMap.Data {
			if v.Username == usernameFromLogin {
				delete(sessions.SessionMap.Data, k)
			}
		}
		s.IsAuthorized = true
		s.Username = usernameFromUserTable
		s.Expiry = time.Now().Add(120 * time.Second)
		cookie := &http.Cookie{
			Name:  "Maryland",
			Value: "1",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// this defaults the current cookie and session to say that no one is logged in and logs the user out.
func Logout(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	sessions.SessionMap.Delete(s)
	http.SetCookie(w, &http.Cookie{Name: "Maryland", Value: "0", Expires: time.Now().Add(365 * 24 * time.Hour), HttpOnly: true})
	http.Redirect(w, r, "/", http.StatusFound)
}

type Info struct {
	Sess       *sessions.Session
	Comments   []comments.CommentFields
	Posts      []posts.PostFields
	LikedPosts []posts.PostFields
	Post       posts.PostFields
	IsAuthor   bool
	Error      string
	UserInfo   users.UserFields
}

func MessageBoard(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		return
	}

	var (
		userFromSession  string
		emailFromSession string
		hashFromSession  string
		iFromSession     string
	)
	row := UserTable.Data.QueryRow("SELECT * from user WHERE username= ?", s.Username)
	switch err := row.Scan(&userFromSession, &emailFromSession, &hashFromSession, &iFromSession); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned! From Session")
	case nil:
		fmt.Println(userFromSession + " Info Found. ")
	default:
		panic(err)
	}
	items := PostsTable.Get(LikesDislikesTable)
	var data Info
	for i, v := range items {
		if v.Author == s.Username {
			items[i].PostAuthor = true
		}
	}
	data = Info{
		Sess:     s,
		Posts:    items,
		UserInfo: users.UserFields{Username: userFromSession, Email: emailFromSession, Image: iFromSession},
	}

	cookie := cookies.FetchCookies(w, r)
	if cookie.Value == "0" {
		t, _ := template.ParseFiles("./templates/homePagewithoutC.html")
		t.Execute(w, data)
	}
	if cookie.Value == "1" {
		c, _ := r.Cookie("sessionId")
		session := sessions.SessionMap.Get(c.Value)
		if !session.IsAuthorized {
			http.SetCookie(w, &http.Cookie{Name: "Maryland", Value: "0", Expires: time.Now().Add(365 * 24 * time.Hour), HttpOnly: true})
			t, _ := template.ParseFiles("./templates/homePagewithoutC.html")
			t.Execute(w, nil)
			return
		}

		//fmt.Println("message board", data)
		t, _ := template.ParseFiles("./templates/homePagewithC.html")
		t.Execute(w, data)
	}
}

func AlreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie(sessions.COOKIE_NAME)
	//fmt.Println(c, "cookies")
	if err != nil {
		return false
	}
	sess := sessions.SessionMap.Data[c.Value]
	if sess.Username != "" {
		return true
	}
	return false
}

func deletePost(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	id := strings.TrimPrefix(r.URL.RequestURI(), "/delete?")
	PostsTable.Delete(id)
	http.Redirect(w, r, "/", http.StatusFound)
}

func savePost(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if s.Username == "" {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.FormValue("content") != "" {
		if r.FormValue("id") != "" {
			id := r.FormValue("id")
			r.ParseForm()
			var thread string
			threadList := r.Form["thread"]
			if len(threadList) > 1 {
				for i, v := range threadList {
					thread += v
					if i != len(threadList)-1 {
						thread += ":"
					}

				}
			} else {
				thread = threadList[0]
			}
			PostsTable.Update(posts.PostFields{
				Content: r.FormValue("content"),
				Thread:  thread,
			}, id)
		} else {
			r.ParseForm()
			var thread string
			threadList := r.Form["thread"]
			if len(threadList) > 1 {
				for i, v := range threadList {
					thread += v
					if i != len(threadList)-1 {
						thread += ":"
					}

				}
			} else {
				thread = threadList[0]
			}
			fmt.Println("about to add 1")
			PostsTable.Add(posts.PostFields{
				Id:      sessions.Generate(),
				Author:  s.Username,
				Content: r.FormValue("content"),
				Thread:  thread,
			})
			fmt.Println("postable", PostsTable)
			fmt.Println(r.FormValue("content"))
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		t, err := template.ParseFiles("./templates/newpost.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data := Info{
			Sess:  s,
			Error: "Silly Billy Fam! You need to write something!!!",
			Post: posts.PostFields{
				Content: r.FormValue("content"),
			},
		}
		t.Execute(w, data)

	}
}

func newPost(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if s.Username == "" {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	items := PostsTable.Get(LikesDislikesTable)
	id := strings.TrimPrefix(r.URL.RequestURI(), "/edit?")
	var item posts.PostFields
	for _, v := range items {
		if v.Id == id {
			item = v
		}
	}
	t, err := template.ParseFiles("./templates/newpost.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data := Info{
		Sess: s,
		Post: item,
	}
	t.Execute(w, data)
}

func LikeDislike(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var value string
	if r.URL.Path == "/like" {
		value = "l"
	} else if r.URL.Path == "/dislike" {
		value = "d"
	}
	LikesDislikesTable.Add(likes.LikesFields{
		PostId:   r.FormValue("id"),
		Username: s.Username,
		Like:     value,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func View(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	var url string
	if r.URL.Path != "/view" && s.Username != "" {
		url = "comment"
	} else {
		url = "view"
	}

	items := PostsTable.Get(LikesDislikesTable)
	// id := strings.TrimPrefix(r.URL.RequestURI(), "/"+url+"?")
	id := r.FormValue("id")
	var item posts.PostFields

	for _, v := range items {
		if v.Id == id {
			item = v
		}
	}
	coms := CommentTable.Get(LikesDislikesCommentsTable, id)
	for i, v := range coms {
		if v.Author == s.Username {
			coms[i].CommentAuthor = true
		}
	}

	data := Info{
		Sess:     s,
		Comments: coms,
		Post:     item,
		IsAuthor: item.Author == s.Username,
	}
	fmt.Println("coms--->",coms)

	t, err := template.ParseFiles("./templates/comment.html", "./templates/view.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, url, data)
}

func LikeDislikecom(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if !s.IsAuthorized {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var value string
	if r.URL.Path == "/likecom" {
		value = "l"
	} else if r.URL.Path == "/dislikecom" {
		value = "d"
	}
	values, _ := url.ParseQuery(r.URL.RawQuery)
	comid := values.Get("coid")
	posid := values.Get("posid")
	LikesDislikesCommentsTable.Add(commentsAndLikes.CommentsAndLikesFields{
		CommentId: comid,
		Username:  s.Username,
		Like:      value,
	})
	http.Redirect(w, r, "/view?id="+posid, http.StatusFound)
}

// func Cabinet(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
// 	myposts, mylikes := PostsTable.GetMyPosts(LikesDislikesTable, s.Username)
// 	data := Info{
// 		Sess:       s,
// 		Posts:      myposts,
// 		LikedPosts: mylikes,
// 	}
// 	t, err := template.ParseFiles("template/cabinet.html", "template/header.html", "template/footer.html")
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Println(mylikes)
// 	t.ExecuteTemplate(w, "cabinet", data)
// }

func Filter(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	thread := r.FormValue("thread")
	if thread == "none" {
		http.Redirect(w, r, "/", 302)
	}
	p := PostsTable.Filter(LikesDislikesTable, thread)
	data := Info{
		Sess:  s,
		Posts: p,
	}
	t, err := template.ParseFiles("./templates/homePagewithoutC.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if s.Username != "" {
		t, err = template.ParseFiles("./templates/homePagewithC.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	t.Execute(w, data)
}

func DelComm(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	values, _ := url.ParseQuery(r.URL.RawQuery)
	comid := values.Get("coid")
	posid := values.Get("posid")
	CommentTable.Delete(comid)
	http.Redirect(w, r, "/view?id="+posid, http.StatusFound)
}

func SaveComm(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if r.FormValue("content") != "" {
		CommentTable.Add(comments.CommentFields{
			CommentId: sessions.Generate(),
			PostId:    r.FormValue("id"),
			Author:    s.Username,
			Content:   r.FormValue("content"),
		})
	}

	http.Redirect(w, r, "/view?id="+r.FormValue("id"), http.StatusFound)
}

// this initialises a test sqlite database and creates a table containing user information.
func initDB() {
	db, _ := sql.Open("sqlite3", "forumDataBase.db")
	UserTable = users.CreateUserTable(db)
	CommentTable = comments.NewCommentTable(db)
	LikesDislikesTable = likes.CreateLikesTable(db)
	LikesDislikesCommentsTable = commentsAndLikes.CreateLikesAndCommentsTable(db)
	PostsTable = posts.CreatePostTable(db)

}

func main() {
	initDB()
	mux := http.NewServeMux()

	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("/dp-images/", http.StripPrefix("/dp-images/", http.FileServer(http.Dir("./dp-images"))))
	mux.HandleFunc("/write", sessions.Middleware(newPost))
	mux.HandleFunc("/edit", sessions.Middleware(newPost))
	mux.HandleFunc("/delete", sessions.Middleware(deletePost))
	mux.HandleFunc("/SavePost", sessions.Middleware(savePost))
	mux.HandleFunc("/like", sessions.Middleware(LikeDislike))
	mux.HandleFunc("/dislike", sessions.Middleware(LikeDislike))
	mux.HandleFunc("/view", sessions.Middleware(View))
	mux.HandleFunc("/comment", sessions.Middleware(View))
	mux.HandleFunc("/savecomm", sessions.Middleware(SaveComm))
	mux.HandleFunc("/deleteCom", sessions.Middleware(DelComm))
	mux.HandleFunc("/filter", sessions.Middleware(Filter))
	// mux.HandleFunc("/cabinet", sessions.Middleware(Cabinet))
	mux.HandleFunc("/likecom", sessions.Middleware(LikeDislikecom))
	mux.HandleFunc("/dislikecom", sessions.Middleware(LikeDislikecom))
	mux.HandleFunc("/", sessions.Middleware(MessageBoard))
	mux.HandleFunc("/login", sessions.Middleware(logIn))
	mux.HandleFunc("/AuthoriseLogin", sessions.Middleware(AuthoriseLogin))
	mux.HandleFunc("/signup", sessions.Middleware(signUp))
	mux.HandleFunc("/avatar", sessions.Middleware(avatar))
	mux.HandleFunc("/upload", sessions.Middleware(uploadFile))
	mux.HandleFunc("/Logout", sessions.Middleware(Logout))
	fmt.Println("Starting Server")
	fmt.Println("Please open http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", mux))
	fmt.Println("error")
}
