package main

import (
	"log"
	"net/http"

	"github.com/guilhermef/go-shortener/handler"
)

const addr = "localhost:12345"

func main() {
	err := http.ListenAndServe(addr, &handler.RedirectHandler{})
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
