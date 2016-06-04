package main

import (
	"log"
	"net/http"
)

const addr = "localhost:12345"

func main() {
	err := http.ListenAndServe(addr, &RedirectHandler{})
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
