package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func init() {
	scenarios = append(scenarios, RunValidScenario)
}

func RunValidScenario(url, email string, nodes []Node, client *http.Client) (errors []error) {
	var buf bytes.Buffer
	resp, err := client.Post(url, "application/json", &buf)
	if err != nil {
		errors = append(errors, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		errors = append(errors, fmt.Errorf("Expected status code 200 when posting valid json, but got %d", resp.StatusCode))
	}
	var validRequestResponse validResponse
	err = json.NewDecoder(resp.Body).Decode(&validRequestResponse)
	if err != nil {
		errors = append(errors, fmt.Errorf("Cannot parse json: %s", err))
	}
	if items := len(validRequestResponse.Response); items != 7 {
		errors = append(errors, fmt.Errorf("Expected %d elements but there were %d", 7, items))
		return
	}
	tests := []struct {
		expected   interface{}
		actual     interface{}
		messageFmt string
	}{
		{"16 Kids and Counting", validRequestResponse.Response[0].Title, "Expected the first item to be '%s', but it was '%s'"},
	}
	for _, test := range tests {
		if test.expected != test.actual {
			errors = append(errors, fmt.Errorf(test.messageFmt, test.expected, test.actual))
		}
	}
	return
}
