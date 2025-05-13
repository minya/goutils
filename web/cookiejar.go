// Package web provides utilities for web requests and HTTP operations.
package web

import (
	"net/http"
	"net/url"
	"sync"
)

// Jar implements http.CookieJar interface for managing HTTP cookies.
// It provides thread-safe storage and retrieval of cookies for HTTP requests.
type Jar struct {
	lk      sync.Mutex
	cookies map[string]map[string]*http.Cookie
}

// NewJar creates and initializes a new cookie jar.
// Returns a pointer to a Jar instance ready to store and retrieve cookies.
func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string]map[string]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL. It stores all cookies associated with a URL host in the jar.
// This method is thread-safe and implements the http.CookieJar interface.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	if jar.cookies[u.Host] == nil {
		jar.cookies[u.Host] = make(map[string]*http.Cookie)
	}

	for _, cookie := range cookies {
		jar.cookies[u.Host][cookie.Name] = cookie
	}
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It returns all cookies stored for the URL's host.
// This method is thread-safe and implements the http.CookieJar interface.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	r := make([]*http.Cookie, len(jar.cookies[u.Host]))
	idx := 0
	for _, c := range jar.cookies[u.Host] {
		r[idx] = c
		idx++
	}
	return r
}
