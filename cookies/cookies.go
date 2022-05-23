package cookies

import (
	"net/http"
	"time"
)

type Cookie struct {
	Name       string
	Value      string
	Path       string    //indicates a URL path that must exist i the req URL to send Cookie Header
	Domain     string    //specifices which hosts are allowed to receive the cookie
	Expires    time.Time //deletes at a specific date
	RawExpires string    //for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int  //deletes cookie after specified amount of time in seconds
	Secure   bool //sent to server on encrypted request
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}

func SetCookies(w http.ResponseWriter, r *http.Request) *http.Cookie {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	//fmt.Println("helllllooooooo")
	cookie := &http.Cookie{Name: "Maryland", Value: "0", Expires: expiration, HttpOnly: true}
	http.SetCookie(w, cookie)
	//fmt.Println("set cookie", cookie)
	return cookie
}

func FetchCookies(w http.ResponseWriter, r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("Maryland")
	//fmt.Println("cookies:", cookie, "err:", err)
	if err != nil {
		//fmt.Println("No Cookies Found")
		cookie = SetCookies(w, r)
	}
	return cookie
}
