package web

import (
	"context"
	"net"
	"net/http"
	"time"
)

// DefaultTimeout is the default timeout duration for HTTP requests
var DefaultTimeout = 10 * time.Second

// dialContextWithTimeout creates a dialer with the specified timeout in milliseconds
func dialContextWithTimeout(ms int) func(ctx context.Context, network, addr string) (net.Conn, error) {
	timeout := time.Duration(ms) * time.Millisecond
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout:   timeout,
			KeepAlive: 30 * time.Second,
		}
		return dialer.DialContext(ctx, network, addr)
	}
}

// DefaultTransport returns net/http's default transport (proxy from the
// HTTPS_PROXY/HTTP_PROXY/NO_PROXY environment, TLS handshake and idle
// timeouts, HTTP/2) with the dial timeout set to ms milliseconds.
func DefaultTransport(ms int) *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DialContext = dialContextWithTimeout(ms)
	return t
}
