package web

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

// Cookies and SetCookies are called concurrently by net/http for requests in
// flight on one client; run with -race.
func TestJarConcurrentUse(t *testing.T) {
	jar := NewJar()
	u, _ := url.Parse("https://example.com/")
	var wg sync.WaitGroup
	for i := range 4 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for j := range 500 {
				jar.SetCookies(u, []*http.Cookie{{Name: fmt.Sprintf("c%d-%d", i, j), Value: "v"}})
			}
		}()
		go func() {
			defer wg.Done()
			for range 500 {
				_ = jar.Cookies(u)
			}
		}()
	}
	wg.Wait()
	if n := len(jar.Cookies(u)); n != 2000 {
		t.Fatalf("jar holds %d cookies, want 2000", n)
	}
}

func TestJarReturnsCookiesPerHost(t *testing.T) {
	jar := NewJar()
	a, _ := url.Parse("https://a.example/")
	b, _ := url.Parse("https://b.example/")
	jar.SetCookies(a, []*http.Cookie{{Name: "s", Value: "1"}})
	jar.SetCookies(a, []*http.Cookie{{Name: "s", Value: "2"}})
	if got := jar.Cookies(a); len(got) != 1 || got[0].Value != "2" {
		t.Fatalf("a cookies = %v", got)
	}
	if got := jar.Cookies(b); len(got) != 0 {
		t.Fatalf("b cookies = %v", got)
	}
}

// F62: deletion, expired and Secure cookies are handled per RFC 6265, and a
// Domain cookie reaches subdomains.
func TestJarFollowsCookieAttributes(t *testing.T) {
	jar := NewJar()
	https, _ := url.Parse("https://rutracker.org/forum/login.php")
	plain, _ := url.Parse("http://rutracker.org/forum/index.php")
	sub, _ := url.Parse("https://bt.rutracker.org/x")

	jar.SetCookies(https, []*http.Cookie{
		{Name: "bb_session", Value: "s", Secure: true},
		{Name: "dead", Value: "d", Expires: time.Unix(1, 0)},
		{Name: "domainwide", Value: "w", Domain: ".rutracker.org", Path: "/"},
	})
	names := func(u *url.URL) map[string]string {
		m := map[string]string{}
		for _, c := range jar.Cookies(u) {
			m[c.Name] = c.Value
		}
		return m
	}

	if got := names(https); got["bb_session"] != "s" || got["domainwide"] != "w" {
		t.Errorf("https cookies = %v, want bb_session and domainwide", got)
	}
	if _, ok := names(https)["dead"]; ok {
		t.Error("an expired cookie is still sent")
	}
	if _, ok := names(plain)["bb_session"]; ok {
		t.Error("a Secure cookie is sent over plain http")
	}
	if got := names(sub); got["domainwide"] != "w" {
		t.Errorf("subdomain cookies = %v, want the Domain cookie", got)
	}

	// Logout: the server deletes the session with Max-Age<0.
	jar.SetCookies(https, []*http.Cookie{{Name: "bb_session", Value: "", MaxAge: -1}})
	if _, ok := names(https)["bb_session"]; ok {
		t.Error("a deleted cookie is still sent")
	}
}
