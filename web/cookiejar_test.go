package web

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"
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
