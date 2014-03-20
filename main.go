package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"sync"
)

var nodes NodeSlice

func errorsAsJson(name, url string, errors []error) string {
	messages := make([]string, len(errors))
	fmt.Printf("Errors found in %s's attempt (%s):", name, url)
	for i, err := range errors {
		println("  " + err.Error())
		messages[i] = err.Error()
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&map[string][]string{"errors": messages})
	return buf.String()
}

func mainHandler(r *http.Request, params martini.Params) (responseCode int, body string) {
	url := r.FormValue("url")
	email := r.FormValue("email")
	name := r.FormValue("name")
	if url == "" || email == "" || name == "" {
		return 400, `{ "error": "'url', 'email' and 'name' must be provided" }`
	}
	client := httpClient(r)
	errors := RunScenarios(url, email, nodes, client)
	if errors != nil {
		SendFailure(name, email, url)
		return 502, errorsAsJson(name, url, errors)
	}
	err := SendConfirmation(name, email, url)
	if err != nil {
		return 500, `{ "error": "` + err.Error() + `" }`
	}
	return 202, `{ "success": true }`
}

var m *martini.Martini

func init() {
	m = martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=60000")
	})
	m.Use(martini.Static("public"))
	m.Use(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})
	r := martini.NewRouter()
	r.Put("/candidates", mainHandler)
	m.Action(r.Handle)
	http.Handle("/", m)
}

// e.g.
// PUT http://localhost:3000/candidates
// email=my@email.com&url=http://my.remote.service
//
// will run scenarios against http://my.remote.service.
//
// If the scenarios pass, should return 202 Accepted with a confirmation response,
// and send us an email with the candidate's details.
//
func main() {
	nodes = NodeSlice{{"8080", true}, {"8081", true}}
	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go func(n Node) {
			http.ListenAndServe(":"+n.Port, &n)
			wg.Done()
		}(node)
	}
	m.Run()
	wg.Wait()
}
