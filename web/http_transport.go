package web

import (
	"net"
	"net/http"
	"time"
)

var timeout = time.Duration(10000 * time.Millisecond)

func dialWithTimeout(secs int) func(network, addr string) (net.Conn, error) {
	return func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, timeout)
	}
}

//proxyUrl, _ := url.Parse("https://192.168.14.140:8888")
func DefaultTransport(ms int) *http.Transport {
	transport := new(http.Transport)
	transport.Dial = dialWithTimeout(ms)
	//transport.Proxy = http.ProxyURL(proxyUrl),
	//transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true},
	return transport
}
