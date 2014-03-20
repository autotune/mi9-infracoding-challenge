package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func init() {
	scenarios = append(scenarios, RunValidScenario)
}

func RunValidScenario(url, email string, nodes []Node, client *http.Client) (errors []error) {
	/*  - Generate a bunch of requests
	*   - Interpret the requests, pulling out the node id
	*   - Count the number of requests per id
	*		- Find the ratio of responses
	 */
	var loadBlancingResults = map[string]int{}
	var validRequestResponse validResponse

	for _, node := range nodes {
		loadBlancingResults[node.Port] = 0
	}

	for i := 0; i < 100; i++ {
		resp, err := client.Get(url)
		if err != nil {
			errors = append(errors, err)
			return
		}
		//defer resp.Body.Close()
		if resp.StatusCode != 200 {
			errors = append(errors, fmt.Errorf("Expected status code 200, but got %d", resp.StatusCode))
		}

		err = json.NewDecoder(resp.Body).Decode(&validRequestResponse)
		if err != nil {
			errors = append(errors, fmt.Errorf("Cannot parse json: %s", err))
		}
		loadBlancingResults[validRequestResponse.Id] = loadBlancingResults[validRequestResponse.Id] + 1
		resp.Body.Close()
	}
	if loadBlancingResults[nodes[0].Port] != loadBlancingResults[nodes[1].Port] {
		errors = append(errors, fmt.Errorf("Load wasn't balanced enough, node0 received %d requests and node1 %d", loadBlancingResults[nodes[0].Port], loadBlancingResults[nodes[1].Port]))
	}
	return
}
