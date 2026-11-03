// Package web provides utilities for web requests and HTTP operations.
package web

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

// Jar implements http.CookieJar following RFC 6265: it honours Expires and
// Max-Age (expired and deletion cookies are dropped), Secure (not sent over
// plain http), Domain (sent to subdomains when set) and Path. It is safe for
// concurrent use. It wraps net/http/cookiejar with the public suffix list,
// so a server cannot set cookies for a whole public suffix such as .co.uk.
type Jar struct {
	inner *cookiejar.Jar
}

// NewJar creates and initializes a new cookie jar.
func NewJar() *Jar {
	inner, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		// cookiejar.New never fails in the standard library.
		panic(err)
	}
	return &Jar{inner: inner}
}

// SetCookies handles the receipt of the cookies in a reply for the given URL.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.inner.SetCookies(u, cookies)
}

// Cookies returns the cookies to send in a request for the given URL.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.inner.Cookies(u)
}
