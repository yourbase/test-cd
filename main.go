package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.Int("port", "8088", "Port for the server to listen")

type helloWorldHandler struct{}

func (h helloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World")
}

func main() {
	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), helloWorldHandler{})
	log.Fatal("HelloWorld ListenAndServe error", err)
}
