package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/codegangsta/martini"
	"sync"
)

type Node struct {
	Port string
	Available bool
}

func (n *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if n.Available {
		fmt.Fprintf(w, "OK from%s - Unique res: %s - %t", n.Port, r.URL.Path[1:], n.Available)
		n.Available = false
	} else {
		http.Error(w, "Backend node failed", 500)
		n.Available = true
	}
}

func errorsAsJson(errors []error) string {
	messages := make([]string, len(errors))
	for i, err := range errors {
		println(err.Error())
		messages[i] = err.Error()
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&map[string][]string{"errors": messages})
	return buf.String()
}

func mainHandler(r *http.Request, params martini.Params) (responseCode int, body string) {
	url := r.FormValue("url")
	email := r.FormValue("email")
	if url == "" || email == "" {
		return 400, `{ "error": "'url' and 'email' must be provided" }`
	}
	errors := RunScenarios(url, email)
	if errors != nil {
		return 502, errorsAsJson(errors)
	}
	return 202, `{ "success": true }`
}

func main() {
	//flag.Parse()
	nodes := []Node{{":8080", true},{":8081", true}}
	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go func(n Node) {
			http.ListenAndServe(n.Port, &n)
			wg.Done()
		}(node)
	}
	m := martini.Classic()
	m.Put("/candidates", mainHandler)
	m.Use(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})
	m.Run()
	wg.Wait()
}
