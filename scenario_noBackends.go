package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func init() {
	scenarios = append(scenarios, RunBadRequestScenario)
}

func RunBadRequestScenario(url, email string, nodes NodeSlice, client *http.Client) (errors []error) {
	nodes.Disable()
	resp, err := client.Get(url)
	if err != nil {
		errors = append(errors, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 503 {
		errors = append(errors, fmt.Errorf("Expected status code 503 when backends are unavailable, but got %d", resp.StatusCode))
	}
	if contentType := resp.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
		errors = append(errors, fmt.Errorf("Expected Content-Type: application/json in response to a bad request, but got %s", contentType))
	}
	var badRequestResponse errorResponse
	err = json.NewDecoder(resp.Body).Decode(&badRequestResponse)
	if err != nil {
		errors = append(errors, err)
	}
	if !strings.Contains(badRequestResponse.Error, "no backend nodes available") {
		errors = append(errors, fmt.Errorf(`The "error" key in the response didn't contain the string 'no backend nodes available'`))
	}
	return
}
