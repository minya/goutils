package web

import (
	"net/http"
	"net/url"
	"testing"
)

// F67: the transport honours HTTPS_PROXY, so the bot can reach Telegram
// and rutracker through a proxy configured in the environment.
func TestDefaultTransportUsesProxyFromEnvironment(t *testing.T) {
	tr := DefaultTransport(1000)
	if tr.Proxy == nil {
		t.Fatal("Proxy is nil: HTTPS_PROXY is ignored")
	}
	// http.ProxyFromEnvironment caches the environment on first use, so
	// check the wiring through the function rather than by setting env.
	req := &http.Request{URL: &url.URL{Scheme: "https", Host: "rutracker.org"}}
	if _, err := tr.Proxy(req); err != nil {
		t.Fatalf("Proxy(req) = %v", err)
	}
	if tr.TLSHandshakeTimeout == 0 || tr.IdleConnTimeout == 0 || tr.DialContext == nil {
		t.Errorf("timeouts: handshake %v, idle %v, dial set %v", tr.TLSHandshakeTimeout, tr.IdleConnTimeout, tr.DialContext != nil)
	}
}
