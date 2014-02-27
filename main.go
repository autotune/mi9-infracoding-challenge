package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Node struct {
	Port string
}

func (n Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK from%s - Unique res: %s", n.Port, r.URL.Path[1:])
}

func main() {
	//flag.Parse()
	ports := []string{":8080", ":8081", ":8082"}
	var wg sync.WaitGroup
	for _, port := range ports {
		wg.Add(1)
		go func(port string) {
			var node = &Node{}
			node.Port = port
			http.ListenAndServe(port, node)
			wg.Done()
		}(port)
	}
	wg.Wait()
}
