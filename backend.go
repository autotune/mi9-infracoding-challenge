package main

import (
	"fmt"
	"net/http"
)

type Node struct {
	Port      string
	Available bool
}

func (n *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if n.Available {
		fmt.Fprintf(w, "{ \"id\": \"%s\",\"response\": \"%s\" }", n.Port, r.URL.Path[1:])
	} else {
		http.Error(w, "Backend node failed", 500)
	}
}
