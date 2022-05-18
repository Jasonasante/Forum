package sessions

import (
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// sessions
const (
	COOKIE_NAME = "sessionId"
)

type Session struct {
	Id           string
	Username     string
	IsAuthorized bool
	Expiry time.Time
}

type SessionStore struct {
	Data map[string]*Session
}

func Generate() string {
	u2 := uuid.NewV4()
	return fmt.Sprintf("%x", u2)
}

var SessionMap = NewSessionStore()

func NewSessionStore() *SessionStore {
	s := new(SessionStore)
	s.Data = make(map[string]*Session)
	return s
}

func (store *SessionStore) Get(sessionId string) *Session {
	session := store.Data[sessionId]
	if session == nil {
		return &Session{Id: sessionId}
	}
	return session
}

func (store *SessionStore) Set(session *Session) {
	store.Data[session.Id] = session
}

func (store *SessionStore) Delete(session *Session) {
	delete(store.Data, session.Id)
}

func IsExpired(s *Session) bool {
	return s.Expiry.Before(time.Now())
}

func ensureSession(r *http.Request, w http.ResponseWriter) string {
	cookie, _ := r.Cookie(COOKIE_NAME)
	if cookie != nil {
		if cookie.Expires.Before(time.Now()) {
			cookie.Expires = time.Now().Add(365 * 24 * time.Hour)
			http.SetCookie(w, cookie)

		}
		return cookie.Value
	}
	sessionId := Generate()
	cookie = &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   sessionId,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}
	http.SetCookie(w, cookie)

	return sessionId
}

func Middleware(next func(w http.ResponseWriter, r *http.Request, s *Session)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := ensureSession(r, w)

		session := SessionMap.Get(sessionId)
		fmt.Println(session, "sessions")

		SessionMap.Set(session)
		next(w, r, session)
	}
}
