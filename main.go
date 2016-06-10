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
	err := http.ListenAndServe("localhost:"+cfg.Port, &handler.RedirectHandler{Client: cfg.RedisClient, Logger: cfg.Logger})
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
