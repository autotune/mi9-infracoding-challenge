package main

import (
	"fmt"
	"net/http"
)

type NodeSlice []*Node

type Node struct {
	Port      string
	Available bool
}

func (n *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if n.Available {
		fmt.Fprintf(w, "{ \"id\": \"%s\",\"Message\": \"%s\" }", n.Port, r.URL.Path[1:])
	} else {
		http.Error(w, "{ \"id\": \""+n.Port+"\",\"status\": \"offline\"}", 503)
	}
}

func (nodes NodeSlice) Enable() {
	for _, node := range nodes {
		node.Available = true
	}
	fmt.Println("Enabling Nodes")
}

func (nodes NodeSlice) Disable() {
	for _, node := range nodes {
		node.Available = false
	}
	fmt.Println("Disabling Nodes")
}
