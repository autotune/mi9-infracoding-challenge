package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func init() {
	scenarios = append(scenarios, RunValidScenario)
}

func TestUrl(url string, client *http.Client) (id string, errors []error) {
	var validRequestResponse validResponse

	resp, err := client.Get(url)
	if err != nil {
		errors = append(errors, err)
		return "", errors
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		errors = append(errors, fmt.Errorf("Expected status code 200, but got %d", resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(&validRequestResponse)
	if err != nil {
		errors = append(errors, fmt.Errorf("Cannot parse json: %s", err))
	}

	return validRequestResponse.Id, errors
}

func RunValidScenario(url, email string, nodes NodeSlice, client *http.Client) (errors []error) {
	/*  - Generate a bunch of requests
	*   - Interpret the requests, pulling out the node id
	*   - Count the number of requests per id
	*		- Find the ratio of responses
	 */
	nodes.Enable()
	time.Sleep(6 * time.Second)
	var loadBlancingResults = map[string]int{}

	for _, node := range nodes {
		loadBlancingResults[node.Port] = 0
	}

	for i := 0; i < 100; i++ {
		id, err := TestUrl(url, client)
		if err != nil {
			errors = append(errors, err...)
		}
		if id != "" {
			loadBlancingResults[id] = loadBlancingResults[id] + 1
		}
	}
	if loadBlancingResults[nodes[0].Port] != loadBlancingResults[nodes[1].Port] {
		errors = append(errors, fmt.Errorf("Load wasn't balanced enough, node0 received %d requests and node1 %d", loadBlancingResults[nodes[0].Port], loadBlancingResults[nodes[1].Port]))
	}
	return
}
