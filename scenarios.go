package main

import (
	"net/http"
)

var scenarios []func(url, email string, nodes NodeSlice, client *http.Client) []error

type errorResponse struct {
	Error string
}

type validResponse struct {
	Id      string
	Message string
}

func RunScenarios(url, email string, nodes NodeSlice, client *http.Client) (errors []error) {
	err := RunAvailabilityScenario(url, client)
	if err != nil {
		errors = append(errors, err)
		return
	}
	for _, scenario := range scenarios {
		errs := scenario(url, email, nodes, client)
		if errs != nil {
			errors = append(errors, errs...)
		}
	}
	return
}

func RunAvailabilityScenario(url string, client *http.Client) error {
	resp, err := client.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
