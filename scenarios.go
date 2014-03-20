package main

import (
	"net"
	"net/http"
	"time"
)

var scenarios []func(url, email string) []error

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, 2*time.Second)
}

func init() {
	// Ensure that the upstream service isn't too slow
	http.DefaultTransport = &http.Transport{
		ResponseHeaderTimeout: 3 * time.Second,
		Dial: dialTimeout,
	}
}

type errorResponse struct {
	Error string
}

type validResponse struct {
	Candidate string
	Response  []struct {
		Image string
		Title string
		Slug  string
	}
}

func RunScenarios(url, email string) (errors []error) {
	err := RunAvailabilityScenario(url)
	if err != nil {
		errors = append(errors, err)
		return
	}
	for _, scenario := range scenarios {
		errs := scenario(url, email)
		if errs != nil {
			errors = append(errors, errs...)
		}
	}
	return
}

func RunAvailabilityScenario(url string) error {
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
