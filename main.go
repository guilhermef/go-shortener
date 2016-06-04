package main

import (
	"github.com/guilhermef/go-shortener/handler"
	//"gopkg.in/redis.v3"
	"log"
	"net/http"
)

const addr = "localhost:12345"

func main() {
	err := http.ListenAndServe(addr, &handler.RedirectHandler{})
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
