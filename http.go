// +build !appengine

package main

import (
	"net"
	"net/http"
	"time"
)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, 2*time.Second)
}

func httpClient(r *http.Request) *http.Client {
	// Ensure that the upstream service isn't too slow
	return &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: 3 * time.Second,
			Dial: dialTimeout,
		},
	}
}
