package utils

import (
	"net"
	"net/http"
	"time"
)

type Client struct {
	*http.Client
}

func HttpClient() *Client {
	const (
		timeout           = time.Second * 10
		handshakeTimeout  = time.Second * 5
		disableKeepAlives = false
	)
	return &Client{
		Client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				DisableKeepAlives: disableKeepAlives,
				Dial: (&net.Dialer{
					Timeout: timeout,
				}).Dial,
				TLSHandshakeTimeout: handshakeTimeout,
			},
		},
	}
}

func (c *Client) Timeout(d time.Duration) *Client {
	c.Client.Timeout = d
	return c
}

func (c *Client) Transport(v http.RoundTripper) *Client {
	c.Client.Transport = v
	return c
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}
