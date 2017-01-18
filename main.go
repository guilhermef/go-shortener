package main

import (
	"github.com/guilhermef/go-shortener/config"
	"github.com/guilhermef/go-shortener/handler"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	address := "0.0.0.0:" + cfg.Port
	log.Printf("Running on %s", address)
	err = http.ListenAndServe(address, &handler.RedirectHandler{Client: cfg.RedisClient, Logger: cfg.Logger })
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
