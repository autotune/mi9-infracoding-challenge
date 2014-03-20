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
		fmt.Fprintf(w, "{ \"id\": \"%s\",\"Message\": \"%s\" }", n.Port, r.URL.Path[1:])
	} else {
		http.Error(w, "{ \"id\": \"%s\",\"status\": \"offline\"}", 500)
	}
}
