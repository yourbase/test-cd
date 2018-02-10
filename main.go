package main

import (
	"fmt"
	"log"
	"net/http"
)

type helloWorldHandler struct{}

func (h helloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World")
}

func main() {
	err := http.ListenAndServe(":8089", helloWorldHandler{})
	log.Fatal("HelloWorld ListenAndServe error", err)
}
