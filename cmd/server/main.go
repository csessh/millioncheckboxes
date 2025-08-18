package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Go backend 🚀")
}

func main() {
	http.HandleFunc("/", helloHandler)

	addr := ":8080"
	log.Printf("Server starting at %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
