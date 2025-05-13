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

// DefaultTransport creates a customized http.Transport with the specified timeout in milliseconds
func DefaultTransport(ms int) *http.Transport {
	return &http.Transport{
		DialContext:           dialContextWithTimeout(ms),
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}
}
